import {
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings,
  MutableDataFrame,
  FieldType,
} from '@grafana/data';

import { getBackendSrv } from '@grafana/runtime';
import _ from 'lodash';
import { Md5 } from 'ts-md5/dist/md5';
import defaults from 'lodash/defaults';

import { MyQuery, MyDataSourceOptions, defaultQuery } from './types';

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

  async query(options: DataQueryRequest<MyQuery>): Promise<DataQueryResponse> {
    const promises = options.targets.map(async target => {
      const query = defaults(target, defaultQuery);
      const response = await this.doRequest({ query: query.queryText });

      /**
       * In this example, the /api/metrics endpoint returns:
       *
       * {
       *   "datapoints": [
       *     {
       *       Time: 1234567891011,
       *       Value: 12.5
       *     },
       *     {
       *     ...
       *   ]
       * }
       */
      const datapoints = response.data.datapoints;

      const timestamps: number[] = [];
      const values: number[] = [];

      for (let i = 0; i < datapoints.length; i++) {
        timestamps.push(datapoints[i].Time);
        values.push(datapoints[i].Value);
      }

      return new MutableDataFrame({
        refId: query.refId,
        fields: [
          { name: 'Time', type: FieldType.time, values: timestamps },
          { name: 'Value', type: FieldType.number, values: values },
        ],
      });
    });

    return Promise.all(promises).then(data => ({ data }));
  }

  async doRequest(params: Record<string, any>) {
    const ttl = 20 * 60 * 1000; // 20 minutes in milliseconds
    const expirationTime = new Date().getTime() + ttl; // in milliseconds
    const basePath = `/${this.account}/data`;
    const baseParams = `dateToken=${expirationTime}&${new URLSearchParams(params).toString()}`;
    const baseToken = `${basePath}?${baseParams}`;
    const token = Md5.hashStr(`${baseToken}${this.apiKey}`);
    console.log(`${baseToken}${this.apiKey}`);
    console.log(token);

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
