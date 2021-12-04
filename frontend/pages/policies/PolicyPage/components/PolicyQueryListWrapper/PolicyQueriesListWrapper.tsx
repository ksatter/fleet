import React from "react";
import { noop } from "lodash";
import paths from "router/paths";

import { ICampaignQueryResult } from "interfaces/campaign";
import { ITeam } from "interfaces/team";
import TableContainer from "components/TableContainer";
import {
  generateTableHeaders,
  generateDataSet,
} from "./PolicyQueriesTableConfig";
// @ts-ignore
import policySvg from "../../../../../../assets/images/no-policy-323x138@2x.png";

const baseClass = "policies-queries-list-wrapper";
const noPolicyQueries = "no-policy-queries";

const TAGGED_TEMPLATES = {
  hostsByTeamRoute: (teamId: number | undefined | null) => {
    return `${teamId ? `/?team_id=${teamId}` : ""}`;
  },
};

interface IPoliciesListWrapperProps {
  policyQueriesList: ICampaignQueryResult[];
  isLoading: boolean;
  resultsTitle?: string;
  canAddOrRemovePolicy?: boolean;
  tableType?: string;
}

const PoliciesListWrapper = ({
  policyQueriesList,
  isLoading,
  resultsTitle,
  canAddOrRemovePolicy,
  tableType,
}: IPoliciesListWrapperProps): JSX.Element => {
  const { MANAGE_HOSTS } = paths;

  const NoPolicyQueries = () => {
    return (
      <div className={`${noPolicyQueries}__inner`}>
        <p>No hosts are online.</p>
      </div>
    );
  };

  return (
    <div
      className={`${baseClass} ${
        canAddOrRemovePolicy ? "" : "hide-selection-column"
      }`}
    >
      <TableContainer
        resultsTitle={resultsTitle || "policies"}
        columns={generateTableHeaders()}
        data={generateDataSet(policyQueriesList)}
        isLoading={isLoading}
        defaultSortHeader={"name"}
        defaultSortDirection={"asc"}
        manualSortBy
        showMarkAllPages={false}
        isAllPagesSelected={false}
        disablePagination
        primarySelectActionButtonVariant="text-icon"
        primarySelectActionButtonIcon="delete"
        primarySelectActionButtonText={"Delete"}
        emptyComponent={NoPolicyQueries}
        onQueryChange={noop}
        disableCount={tableType === "inheritedPolicies"}
      />
    </div>
  );
};

export default PoliciesListWrapper;
