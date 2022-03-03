import { defaults } from 'lodash';

import React, { PureComponent } from 'react';
import { InlineFormLabel, Select } from '@grafana/ui';
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

const typeOptions: Array<SelectableValue<string>> = [
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

  onTypeChange = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, type: option.value });
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);

    return (
      <>
        <div className="gf-form-inline">
          <InlineFormLabel width={12}>&nbsp;</InlineFormLabel>
          <InlineFormLabel width={5}>Filter</InlineFormLabel>
          <Select
            isSearchable={true}
            width={33}
            value={query.filter || ''}
            options={filterOptions}
            onChange={this.onFilterChange}
          />
          <InlineFormLabel width={5}>Metrics</InlineFormLabel>
          <Select
            isSearchable={true}
            width={33}
            value={query.metrics || ''}
            options={metricsOptions}
            onChange={this.onMetricsChange}
          />
          <InlineFormLabel width={5}>Type</InlineFormLabel>
          <Select
            isSearchable={true}
            width={33}
            value={query.metrics || ''}
            options={typeOptions}
            onChange={this.onTypeChange}
          />
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
      </>
    );
  }
}
