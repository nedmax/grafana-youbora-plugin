import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  filter?: string;
  metrics: string;
  type: string;
}

export const defaultQuery: Partial<MyQuery> = {
  filter: '',
  metrics: 'views',
  type: 'ALL',
};

export const MyParams = ['filter', 'metrics', 'type'];

export type StreamingType = 'ALL' | 'VOD' | 'LIVE';

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
