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
  { label: 'Plays', value: 'views' },
  { label: 'Unique Views', value: 'uniques' },
];

const streamingTypeOptions: Array<SelectableValue<string>> = [
  { label: 'All', value: 'ALL' },
  { label: 'VoD', value: 'VOD' },
  { label: 'Live', value: 'LIVE' },
];

export class QueryEditor extends PureComponent<Props> {
  onFilterChange = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, filter: option.value });
    onRunQuery();
  };

  onMetricsChange = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, metrics: option.value });
    onRunQuery();
  };

  onTypeChange = (options: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;

    console.log('AQUI');
    console.log(options);
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
            onChange={this.onFilterChange}
          />
          <InlineFormLabel width={6}>Metrics</InlineFormLabel>
          <Select
            isSearchable={true}
            width={33}
            value={query.metrics || ''}
            options={metricsOptions}
            onChange={this.onMetricsChange}
          />
          <InlineField label="Streaming Type" labelWidth={14}>
            <MultiSelect
              isSearchable
              isClearable
              value={query.streamingType}
              options={streamingTypeOptions}
              onChange={this.onTypeChange}
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
