import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { InlineFormLabel, InlineField, Select, MultiSelect } from '@grafana/ui';
import { SelectableValue, QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

const filterOptions: Array<SelectableValue<string>> = [
  { label: 'start_date', value: 'start_date' },
  { label: 'name', value: 'name' },
  { label: 'distance', value: 'distance' },
  { label: 'moving_time', value: 'moving_time' },
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
  onFilterChanged = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, filter: option.value });
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
          <InlineFormLabel width={5}>Filter</InlineFormLabel>
          <Select
            isSearchable={true}
            width={33}
            value={query.filter || ''}
            options={filterOptions}
            onChange={this.onFilterChanged}
          />

          <InlineField label="Metrics" labelWidth={14}>
            <MultiSelect
              isSearchable
              isClearable
              value={query.metrics}
              options={metricsOptions}
              onChange={this.onMetricsChanged}
            />
          </InlineField>

          <InlineField label="Streaming Type" labelWidth={14}>
            <MultiSelect
              isSearchable
              isClearable
              value={query.streamingType}
              options={streamingTypeOptions}
              onChange={this.onTypeChanged}
            />
          </InlineField>
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
      </>
    );
  }
}
