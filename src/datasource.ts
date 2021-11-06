import {
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings,
  MutableDataFrame,
  FieldType,
} from '@grafana/data';

import { getBackendSrv } from '@grafana/runtime';
import _, { keys } from 'lodash';
import { Md5 } from 'ts-md5/dist/md5';
import defaults from 'lodash/defaults';

import { MyQuery, MyDataSourceOptions, defaultQuery, MyParams } from './types';

export class DataSource extends DataSourceApi<MyQuery, MyDataSourceOptions> {
  url?: string;
  account: string;
  apiKey: string;

  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);

    this.url = instanceSettings.url;
    this.account = instanceSettings.jsonData.account!;
    this.apiKey = instanceSettings.jsonData.apiKey!;
  }

  async query(request: DataQueryRequest<MyQuery>): Promise<DataQueryResponse> {
    const promises = request.targets.map(async target => {
      const { range } = request;
      const query = defaults(target, defaultQuery);
      query.fromDate = range!.from.valueOf().toString();
      query.toDate = range!.to.valueOf().toString();
      const response = await this.doRequest(query);

      const datapoints = response.data.data;

      const timestamps: number[] = datapoints[0].metrics[0].values[0].data.map((x: number[]) => x[0]);
      const dataframe = new MutableDataFrame({
        refId: query.refId,
        fields: [{ name: 'Time', type: FieldType.time, values: timestamps }],
      });

      for (let i = 0; i < datapoints.length; i++) {
        const values: number[] = datapoints[i].metrics[0].values[0].data.map((x: number[]) => x[1]);
        const name = `${datapoints[i].type} ${datapoints[i].metrics[0].label}`;
        dataframe.addField({ name: name, type: FieldType.number, values: values });
      }

      return dataframe;
    });

    return Promise.all(promises).then(data => ({ data }));
  }

  async doRequest(params: Record<string, any>) {
    const ttl = 20 * 60 * 1000; // 20 minutes in milliseconds
    const expirationTime = new Date().getTime() + ttl; // in milliseconds
    const basePath = `/${this.account}/data`;
    let orderedParams = '';
    keys(params)
      .sort()
      .filter(e => MyParams.includes(e))
      .forEach(e => (orderedParams += `&${e}=${params[e]}`));

    const baseParams = `dateToken=${expirationTime}${orderedParams}`;
    const baseToken = `${basePath}?${baseParams}`;
    const token = Md5.hashStr(`${baseToken}${this.apiKey}`);

    return getBackendSrv().datasourceRequest({
      url: `${this.url}/youbora${basePath}?${baseParams}&token=${token}`,
      method: 'GET',
    });
  }

  /**
   * Checks whether we can connect to the API.
   */
  async testDatasource() {
    const defaultErrorMessage = 'Error accessing API';

    try {
      const response = await this.doRequest({ fromDate: 'last5minutes', metrics: 'views' });
      if (response.status === 200) {
        return {
          status: 'success',
          message: 'Success',
        };
      } else {
        return {
          status: 'error',
          message: response.statusText ? response.statusText : defaultErrorMessage,
        };
      }
    } catch (err) {
      if (_.isString(err)) {
        return {
          status: 'error',
          message: err,
        };
      } else {
        let message = '';
        message += err.statusText ? err.statusText : defaultErrorMessage;
        if (err.data && err.data.error && err.data.error.code) {
          message += ': ' + err.data.error.code + '. ' + err.data.error.message;
        }

        return {
          status: 'error',
          message,
        };
      }
    }
  }
}
