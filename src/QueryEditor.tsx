import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { InlineFormLabel, InlineField, Select, MultiSelect, Input } from '@grafana/ui';
import { SelectableValue, QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

const dimensionOptions: Array<SelectableValue<string>> = [
  // Geo
  { label: 'City', value: 'city' },
  { label: 'Country', value: 'country' },
  { label: 'Region', value: 'region' },
  { label: 'State Province', value: 'state_province' },
  // Network
  { label: 'ISP', value: 'isp' },
  { label: 'ASN', value: 'asn' },
  { label: 'Connection Type', value: 'connectiontype' },
  { label: 'CDN', value: 'cdn' },
  // Device
  { label: 'Device Type', value: 'device_type' },
  { label: 'Device Vendor', value: 'device_vendor' },
  // Error
  { label: 'Error Name', value: 'error_name' },
  { label: 'Error Description', value: 'error_description' },
  { label: 'Playback Status', value: 'view_status' },
  // Extraparams
  // TODO: move to datasource configuration
  { label: 'Media ID', value: 'extraparam1' },
  { label: 'Program ID', value: 'extraparam2' },
  { label: 'Channel ID', value: 'extraparam3' },
  { label: 'Player Version', value: 'extraparam4' },
  { label: 'JS API Version', value: 'extraparam5' },
  { label: 'Kind', value: 'extraparam6' },
  { label: 'Domain', value: 'extraparam7' },
  { label: 'Player Type (Platform)', value: 'extraparam8' },
  { label: 'Experiment Alternative', value: 'extraparam9' },
  { label: 'Playback Name', value: 'extraparam10' },
  { label: 'Signal', value: 'extraparam11' },
  { label: 'AutoPlay', value: 'extraparam12' },
  { label: 'StartAt', value: 'extraparam13' },
  { label: 'Total Attempts Video Session', value: 'extraparam14' },
  { label: 'POP', value: 'extraparam15' },
  { label: 'Fallback Video Session', value: 'extraparam16' },
  { label: 'DAI Played', value: 'extraparam17' },
  { label: 'DAI Session ID', value: 'extraparam18' },
  { label: 'Disable Ads', value: 'extraparam19' },
  { label: 'DAI Enabled', value: 'extraparam20' },
];

const metricsOptions: Array<SelectableValue<string>> = [
  // Audience
  { label: 'Plays', value: 'views' },
  { label: 'Plays Initiated', value: 'successfulPlays' },
  { label: 'Concurrent Plays', value: 'concurrent' },
  { label: 'Active Plays', value: 'activeSessions' },
  { label: 'Unique Views', value: 'uniques' },
  { label: 'Subscribers', value: 'subscribers' },
  { label: 'Effective Playtime', value: 'effectiveHours' },
  { label: 'Avg. Effective Playtime', value: 'effectiveTime' },
  { label: 'Traffic', value: 'traffic' },
  { label: 'Avg. Playtime', value: 'playtime' },
  { label: 'EBVS', value: 'exitsCount' },
  { label: 'EBVS Ratio', value: 'exits' },
  { label: 'EBVS excluding ad drops', value: 'ebvsNoAdsCount' },
  { label: 'EBVS excluding ad drops Ratio', value: 'ebvsNoAds' },
  // Quality
  { label: 'Join Time', value: 'jointime' },
  { label: 'Effective Buffer Ratio', value: 'effectiveBufferRatio' },
  { label: 'Avg. Bitrate (Mbps)', value: 'bitratembps' },
  // Errors
  { label: 'Startup Error (#) (globo)', value: 'startupErrorCount_globo' },
  { label: 'Startup Error (%) (globo)', value: 'startupError_globo' },
  { label: 'In-Stream Error (globo)', value: 'inStreamError_globo' },
];

const streamingTypeOptions: Array<SelectableValue<string>> = [
  { label: 'All', value: 'ALL' },
  { label: 'VoD', value: 'VOD' },
  { label: 'Live', value: 'LIVE' },
];

export class QueryEditor extends PureComponent<Props> {
  onFilterByChanged = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, filterBy: option.value });
    onRunQuery();
  };

  onFilterValueChanged = (option: React.FormEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, filterValue: option.currentTarget.value });
    onRunQuery();
  };

  onGroupByChanged = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, groupBy: option.value });
    onRunQuery();
  };

  onMetricsChanged = (options: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;

    if (options) {
      const values: string[] = [];
      options.forEach((option: { value: string }) => option.value && values.push(option.value));
      onChange({ ...query, metrics: values });
      onRunQuery();
    }
  };

  onTypeChanged = (options: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;

    if (options) {
      const values: string[] = [];
      options.forEach((option: { value: string }) => option.value && values.push(option.value));
      onChange({ ...query, streamingType: values });
      onRunQuery();
    }
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);

    return (
      <>
        <div className="gf-form-inline">
          <InlineFormLabel width={10}>Filter by</InlineFormLabel>
          <Select
            isSearchable={true}
            width={33}
            value={query.filterBy || ''}
            options={dimensionOptions}
            onChange={this.onFilterByChanged}
          />
          <InlineField label="value">
            <Input type="text" value={query.filterValue || ''} onChange={this.onFilterValueChanged} />
          </InlineField>
        </div>
        <div className="gf-form-inline">
          <InlineFormLabel width={10}>Group by</InlineFormLabel>
          <Select
            isSearchable={true}
            width={33}
            value={query.groupBy || ''}
            options={dimensionOptions}
            onChange={this.onGroupByChanged}
          />
        </div>
        <div className="gf-form-inline">
          <InlineFormLabel width={10}>Metrics</InlineFormLabel>
          <MultiSelect
            isSearchable
            isClearable
            width={33}
            value={query.metrics}
            options={metricsOptions}
            onChange={this.onMetricsChanged}
          />
        </div>
        <div className="gf-form-inline">
          <InlineFormLabel width={10}>Streaming Type</InlineFormLabel>
          <MultiSelect
            isSearchable
            width={33}
            isClearable
            value={query.streamingType}
            options={streamingTypeOptions}
            onChange={this.onTypeChanged}
          />
        </div>
      </>
    );
  }
}
