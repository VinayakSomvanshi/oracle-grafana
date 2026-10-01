import React, { useState } from 'react';
import { Button, CodeEditor, Label, TextArea } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';

import { DataSource } from '../datasource';
import { interpolate } from '../interpolate';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ onChange, query }: Props) {
  const [showParsed, setShowParsed] = useState(false);

  const onSQLChange = (value: string) => {
    onChange({
      ...query,
      o_sql: value
    });
  };

  const insertMacro = (macro: string) => {
    const current = query.o_sql || '';
    onSQLChange(current ? `${current} ${macro}` : macro);
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', width: '100%' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Label description="Write raw Oracle SQL query. Supports standard macros ($__timeFilter, $__timeFrom, etc.)">
          Oracle SQL Query
        </Label>
        <div style={{ display: 'flex', gap: '6px' }}>
          <Button
            size="sm"
            variant="secondary"
            fill="outline"
            onClick={() => setShowParsed(!showParsed)}
          >
            {showParsed ? 'Hide Parsed Query' : 'Preview Parsed Query'}
          </Button>
        </div>
      </div>

      <div style={{ border: '1px solid rgba(204, 204, 220, 0.15)', borderRadius: '3px', overflow: 'hidden' }}>
        <CodeEditor
          language="sql"
          value={query.o_sql || ''}
          onBlur={onSQLChange}
          onChange={onSQLChange}
          showLineNumbers={true}
          showMiniMap={false}
          height="240px"
          monacoOptions={{
            fontSize: 13,
            lineNumbersMinChars: 3,
            scrollBeyondLastLine: false,
            wordWrap: 'on'
          }}
        />
      </div>

      <div style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '11px', color: 'rgba(204, 204, 220, 0.65)' }}>
        <span>Quick Macros:</span>
        <button
          type="button"
          onClick={() => insertMacro('$__timeFilter(ts)')}
          style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)', borderRadius: '3px', color: '#79b8ff', cursor: 'pointer', padding: '1px 6px', fontSize: '11px' }}
        >
          $__timeFilter(ts)
        </button>
        <button
          type="button"
          onClick={() => insertMacro('$__timeFrom()')}
          style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)', borderRadius: '3px', color: '#79b8ff', cursor: 'pointer', padding: '1px 6px', fontSize: '11px' }}
        >
          $__timeFrom()
        </button>
        <button
          type="button"
          onClick={() => insertMacro('$__timeTo()')}
          style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)', borderRadius: '3px', color: '#79b8ff', cursor: 'pointer', padding: '1px 6px', fontSize: '11px' }}
        >
          $__timeTo()
        </button>
        <button
          type="button"
          onClick={() => insertMacro('$__interval_ms')}
          style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)', borderRadius: '3px', color: '#79b8ff', cursor: 'pointer', padding: '1px 6px', fontSize: '11px' }}
        >
          $__interval_ms
        </button>
      </div>

      {showParsed && (
        <div style={{ marginTop: '8px' }}>
          <Label description="Executed SQL statement after client-side variable interpolation">
            Parsed Query Preview
          </Label>
          <TextArea
            readOnly
            rows={6}
            value={interpolate(query.o_sql ?? '')}
            width="100%"
          />
        </div>
      )}
    </div>
  );
}
