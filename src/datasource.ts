import { DataSourceInstanceSettings, CoreApp, MetricFindValue, dateTime, DataFrame, DataQueryRequest, DataQueryResponse } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { Observable, lastValueFrom, map, switchMap } from 'rxjs';

import { interpolate } from './interpolate';
import { MyQuery, MyDataSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  query(request: DataQueryRequest<MyQuery>): Observable<DataQueryResponse> {
    for(const query of request.targets) {
      if (request.scopedVars && Object.keys(request.scopedVars).length > 0) {
        query.o_parsed = interpolate(query.o_sql || '', request.scopedVars);
      }
    }
    return super.query(request)
  }

  /**
   * Method implemented to use the Query variable available to this datasource.
   * @param query User defined query.
   * @param options Query options.
   * @returns 
   */
  async metricFindQuery(query: string, options?: any): Promise<MetricFindValue[]> {
    if (!query) {
      return Promise.resolve([]);
    }

    const response = this.query({
      interval: '',
      intervalMs: 0,
      requestId: 'metricFindQuery',
      range: {
        from: dateTime(),
        to: dateTime(),
        raw: {
          from: dateTime(),
          to: dateTime()
        }
      },
      scopedVars: {},
      targets: [{
        datasource: this.getDefaultQuery(CoreApp.Unknown).datasource,
        o_parsed: query,
        refId: 'A'
      }],
      timezone: 'Z',
      app: '',
      startTime: 0,
    });

    return lastValueFrom(response.pipe(
      map(response => {
        const results: MetricFindValue[] = [];
        if (!response.data || response.data.length === 0) {
          return results;
        }

        for (const frame of response.data as DataFrame[]) {
          if (!frame.fields || frame.fields.length === 0) {
            continue;
          }

          let textField = frame.fields[0];
          let valueField = frame.fields[0];

          // Check for standard Grafana __text and __value conventions (case-insensitive for Oracle)
          for (const field of frame.fields) {
            const lower = field.name.toLowerCase();
            if (lower === '__text' || lower === 'text') {
              textField = field;
            } else if (lower === '__value' || lower === 'value') {
              valueField = field;
            }
          }

          // If no explicit __text/__value found and we have >= 2 columns, column 0 is text, column 1 is value
          if (textField === valueField && frame.fields.length >= 2) {
            textField = frame.fields[0];
            valueField = frame.fields[1];
          }

          const rowCount = frame.length;
          for (let i = 0; i < rowCount; i++) {
            const textVal = textField.values.get(i);
            const valVal = valueField.values.get(i);
            results.push({
              text: textVal !== null && textVal !== undefined ? String(textVal) : '',
              value: valVal !== null && valVal !== undefined ? String(valVal) : undefined,
            });
          }
        }
        return results;
      })
    ));
  }

  getDefaultQuery(_: CoreApp): Partial<MyQuery> {
    return {
      o_sql: 'SELECT SYSDATE AS time, 100 AS value FROM DUAL',
    }
  }
}
