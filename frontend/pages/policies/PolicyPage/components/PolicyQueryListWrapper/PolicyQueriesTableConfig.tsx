/* eslint-disable react/prop-types */
// disable this rule as it was throwing an error in Header and Cell component
// definitions for the selection row for some reason when we dont really need it.
import React from "react";
import { memoize } from "lodash";

// @ts-ignore
import TextCell from "components/TableContainer/DataTable/TextCell/TextCell";
import { ICampaignQueryResult } from "interfaces/campaign";
import sortUtils from "utilities/sort";
import PassIcon from "../../../../../../assets/images/icon-check-circle-green-16x16@2x.png";
import FailIcon from "../../../../../../assets/images/icon-exclamation-circle-red-16x16@2x.png";

// TODO functions for paths math e.g., path={PATHS.MANAGE_HOSTS + getParams(cellProps.row.original)}

interface IHeaderProps {
  column: {
    host: string;
    isSortedDesc: boolean;
  };
}

interface ICellProps {
  cell: {
    value: any;
  };
  row: {
    original: ICampaignQueryResult;
  };
}

interface IDataColumn {
  Header: ((props: IHeaderProps) => JSX.Element) | string;
  Cell: (props: ICellProps) => JSX.Element;
  title?: string;
  accessor?: string;
  disableHidden?: boolean;
  disableSortBy?: boolean;
  sortType?: string;
}

// NOTE: cellProps come from react-table
// more info here https://react-table.tanstack.com/docs/api/useTable#cell-properties
const generateTableHeaders = (): IDataColumn[] => {
  const tableHeaders: IDataColumn[] = [
    {
      title: "Host",
      Header: "Host",
      disableSortBy: true,
      accessor: "name",
      Cell: (cellProps: ICellProps): JSX.Element => (
        <TextCell value={cellProps.row.original.host_hostname} />
      ),
    },
    {
      title: "Status",
      Header: "Status",
      disableSortBy: true,
      accessor: "status",
      Cell: (cellProps: ICellProps): JSX.Element => (
        <>
          <img alt="host passing" src={PassIcon} />
          <img alt="host failing" src={FailIcon} />
          <TextCell value={cellProps.cell.value} />
        </>
      ),
    },
  ];
  return tableHeaders;
};

const generateDataSet = memoize(
  (policyQueriesList: ICampaignQueryResult[] = []): ICampaignQueryResult[] => {
    policyQueriesList = policyQueriesList.sort((a, b) =>
      sortUtils.caseInsensitiveAsc(a.host_hostname, b.host_hostname)
    );
    return policyQueriesList;
  }
);

export { generateTableHeaders, generateDataSet };
