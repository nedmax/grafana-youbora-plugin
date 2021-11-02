import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  queryText?: string;
}

export const defaultQuery: Partial<MyQuery> = {
  queryText: 'type=ALL&fromDate=last5minutes',
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  account?: string;
  apiKey?: string;
}
