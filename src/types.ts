import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  filterBy?: string;
  filterValue?: string;
  groupBy?: string;
  metrics: string[];
  streamingType: string[];
}

export const defaultQuery: Partial<MyQuery> = {
  filterBy: '',
  filterValue: '',
  groupBy: '',
  metrics: ['views'],
  streamingType: ['ALL'],
};

export const MyParams = ['filter', 'metrics', 'streamingType'];

/**
 * These are options configured for each DataSource instance.
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  account?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apikey?: string;
}
