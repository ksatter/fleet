import React from "react";
import { Row } from "react-table";

import { IHost } from "interfaces/host";
import { ITarget } from "interfaces/target";
import { filter, isEmpty, pullAllBy } from "lodash";

import DataError from "components/DataError";
// @ts-ignore
import Input from "components/forms/fields/InputFieldWithIcon";
import TableContainer from "components/TableContainer";
import { generateTableHeaders } from "./TargetsInputHostsTableConfig";

interface ITargetsInputProps {
  tabIndex: number;
  searchText: string;
  relatedHosts: IHost[];
  isTargetsLoading: boolean;
  hasFetchError: boolean;
  selectedTargets: ITarget[];
  setSearchText: (value: string) => void;
  handleRowSelect: (value: Row) => void;
  handleRowRemove: (value: Row) => void;
}

const baseClass = "targets-input";

const TargetsInput = ({
  tabIndex,
  searchText,
  relatedHosts,
  isTargetsLoading,
  hasFetchError,
  selectedTargets,
  handleRowSelect,
  handleRowRemove,
  setSearchText,
}: ITargetsInputProps): JSX.Element => {
  const resultsDropdownTableHeaders = generateTableHeaders(false);
  const selectedTableHeaders = generateTableHeaders(true);
  const finalRelatedHosts =
    relatedHosts && pullAllBy(relatedHosts, selectedTargets, "hostname");
  const finalSelectedHostTargets =
    selectedTargets && filter(selectedTargets, "hostname");
  const isActiveSearch =
    !isEmpty(searchText) && (!hasFetchError || isTargetsLoading);
  const isSearchError = !isEmpty(searchText) && hasFetchError;

  return (
    <div className={baseClass}>
      <Input
        autofocus
        type="search"
        iconName="search"
        value={searchText}
        tabIndex={tabIndex}
        iconPosition="start"
        label="Target specific hosts"
        placeholder="Search hosts by hostname, UUID, MAC address, or private IP address"
        onChange={setSearchText}
      />
      {isActiveSearch && (
        <div className={`${baseClass}__hosts-search-dropdown`}>
          <TableContainer
            columns={resultsDropdownTableHeaders}
            data={finalRelatedHosts}
            isLoading={isTargetsLoading}
            resultsTitle=""
            emptyComponent={() => (
              <div className="empty-search">
                <div className="empty-search__inner">
                  <h4>No hosts match the current search criteria.</h4>
                  <p>
                    Expecting to see hosts? Try again in a few seconds as the
                    system catches up.
                  </p>
                </div>
              </div>
            )}
            showMarkAllPages={false}
            isAllPagesSelected={false}
            disableCount
            disablePagination
            disableMultiRowSelect
            onSelectSingleRow={handleRowSelect}
          />
        </div>
      )}
      {isSearchError && (
        <div className={`${baseClass}__hosts-search-dropdown`}>
          <DataError />
        </div>
      )}
      <div className={`${baseClass}__hosts-selected-table`}>
        <TableContainer
          columns={selectedTableHeaders}
          data={finalSelectedHostTargets}
          isLoading={false}
          resultsTitle=""
          showMarkAllPages={false}
          isAllPagesSelected={false}
          disableCount
          disablePagination
          disableMultiRowSelect
          emptyComponent={() => <></>}
          onSelectSingleRow={handleRowRemove}
        />
      </div>
    </div>
  );
};

export default TargetsInput;
