import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  filter?: string;
  timezone: string;
  metrics: string;
  fromDate: string;
  toDate?: string;
  granularity: string;
  type: string;
}

export const defaultQuery: Partial<MyQuery> = {
  filter: '',
  timezone: 'GMT',
  granularity: 'minute',
  metrics: 'views',
  fromDate: 'last6hours',
};

export const MyParams = ['filter', 'timezone', 'granularity', 'fromDate', 'metrics', 'toDate'];

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  account?: string;
  apiKey?: string;
}
