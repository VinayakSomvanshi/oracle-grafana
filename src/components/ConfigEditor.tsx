import React, { ChangeEvent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { Checkbox, HorizontalGroup, InlineField, Input, Legend, SecretInput, VerticalGroup } from '@grafana/ui';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> { }

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;

  const onConnStrChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        o_connStr: event.target.value,
      }
    });
  };

  const onHostnameChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        o_hostname: event.target.value,
      }
    });
  };

  const onPortChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        o_port: Number(event.target.value),
      }
    });
  };

  const onServiceChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        o_service: event.target.value,
        o_sid: ''
      }
    });
  };

  const onSIDChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        o_service: '',
        o_sid: event.target.value,
      }
    });
  };

  const onUserChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        o_user: event.target.value,
      }
    });
  };

  // Secure password field
  const onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        o_password: event.target.value,
      },
    });
  };

  const onPasswordReset = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        o_password: false
      },
      secureJsonData: {
        ...options.secureJsonData,
        o_password: ''
      },
    });
  };

  // TLS & Wallet handlers
  const onTLSChange = (event: ChangeEvent<HTMLInputElement>) => {
    const isChecked = event.target.checked;
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        o_tls: isChecked,
        o_tlsVerify: isChecked ? (options.jsonData.o_tlsVerify ?? true) : false,
      }
    });
  };

  const onTLSVerifyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        o_tlsVerify: event.target.checked,
      }
    });
  };

  const onWalletPathChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        o_walletPath: event.target.value,
      }
    });
  };

  const onWalletPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        o_walletPassword: event.target.value,
      },
    });
  };

  const onWalletPasswordReset = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        o_walletPassword: false
      },
      secureJsonData: {
        ...options.secureJsonData,
        o_walletPassword: ''
      },
    });
  };

  const { jsonData, secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

  return (
    <VerticalGroup>
      <Legend>
        Authentication
      </Legend>
      <HorizontalGroup>
        <InlineField grow label="User" labelWidth={12}>
          <Input
            placeholder="oracle_user"
            required
            value={jsonData.o_user || ''}
            width={40}
            onChange={onUserChange}
          />
        </InlineField>
        <InlineField grow label="Password" labelWidth={12}>
          <SecretInput
            isConfigured={(secureJsonFields && secureJsonFields.o_password) as boolean}
            placeholder="oracle_password"
            required
            value={secureJsonData.o_password || ''}
            width={40}
            onChange={onPasswordChange}
            onReset={onPasswordReset}
          />
        </InlineField>
      </HorizontalGroup>

      <Legend>
        Connection
      </Legend>
      <InlineField grow label="ConnString" labelWidth={12} tooltip="Optional JDBC connection string or TNS connect descriptor. Overrides Hostname/Port if specified.">
        <Input
          placeholder="(DESCRIPTION=(ADDRESS=(PROTOCOL=tcp)(HOST=localhost)(PORT=1521))(CONNECT_DATA=(SID=XE)))"
          value={jsonData.o_connStr || ''}
          width={94}
          onChange={onConnStrChange}
        />
      </InlineField>
      <HorizontalGroup>
        <InlineField grow label="Hostname" labelWidth={12}>
          <Input
            placeholder="localhost"
            value={jsonData.o_hostname || ''}
            width={40}
            onChange={onHostnameChange}
          />
        </InlineField>
        <InlineField grow label="Port" labelWidth={12}>
          <Input
            placeholder="1521"
            type="number"
            value={jsonData.o_port || 1521}
            width={40}
            onChange={onPortChange}
          />
        </InlineField>
      </HorizontalGroup>
      <HorizontalGroup>
        <InlineField grow label="Service" labelWidth={12}>
          <Input
            placeholder="e.g. FREEPDB1 or ORCL"
            value={jsonData.o_service || ''}
            width={40}
            onChange={onServiceChange}
          />
        </InlineField>
        <InlineField grow label="or SID" labelWidth={12}>
          <Input
            placeholder="e.g. XE"
            value={jsonData.o_sid || ''}
            width={40}
            onChange={onSIDChange}
          />
        </InlineField>
      </HorizontalGroup>

      <Legend>
        Security & Oracle Cloud (OCI)
      </Legend>
      <HorizontalGroup>
        <InlineField label="Enable TLS / SSL" labelWidth={16} tooltip="Encrypt network traffic using TLS (TCPS). Required for Oracle Autonomous Database on OCI.">
          <Checkbox
            value={jsonData.o_tls || false}
            onChange={onTLSChange}
          />
        </InlineField>
        {jsonData.o_tls && (
          <InlineField label="Verify Certificate" labelWidth={16} tooltip="Validate server TLS certificate against trusted CAs. Disable for self-signed development certificates.">
            <Checkbox
              value={jsonData.o_tlsVerify ?? true}
              onChange={onTLSVerifyChange}
            />
          </InlineField>
        )}
      </HorizontalGroup>
      <HorizontalGroup>
        <InlineField grow label="Wallet Path" labelWidth={16} tooltip="Path on the Grafana server to directory containing cwallet.sso or ewallet.p12 (for OCI ATP/ADW).">
          <Input
            placeholder="e.g. /etc/oracle/wallets/adw"
            value={jsonData.o_walletPath || ''}
            width={40}
            onChange={onWalletPathChange}
          />
        </InlineField>
        <InlineField grow label="Wallet Password" labelWidth={16} tooltip="Optional password to decrypt ewallet.p12 PKCS#12 wallet files.">
          <SecretInput
            isConfigured={(secureJsonFields && secureJsonFields.o_walletPassword) as boolean}
            placeholder="wallet_password"
            value={secureJsonData.o_walletPassword || ''}
            width={40}
            onChange={onWalletPasswordChange}
            onReset={onWalletPasswordReset}
          />
        </InlineField>
      </HorizontalGroup>
    </VerticalGroup>
  );
}
