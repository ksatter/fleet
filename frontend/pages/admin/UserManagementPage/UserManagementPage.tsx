import React, { useState, useCallback, useContext, useEffect } from "react";
import { InjectedRouter } from "react-router";
import { useQuery } from "react-query";
import memoize from "memoize-one";

import paths from "router/paths";
import { IApiError } from "interfaces/errors";
import { IInvite } from "interfaces/invite";
import { IUser, IUserFormErrors } from "interfaces/user";
import { ITeam } from "interfaces/team";
import { clearToken } from "utilities/local";

import { AppContext } from "context/app";
import { NotificationContext } from "context/notification";
import teamsAPI from "services/entities/teams";
import usersAPI from "services/entities/users";
import invitesAPI from "services/entities/invites";

import TableContainer, { ITableQueryData } from "components/TableContainer";
import TableDataError from "components/DataError";
import Modal from "components/Modal";
import { DEFAULT_CREATE_USER_ERRORS } from "utilities/constants";
import EmptyUsers from "./components/EmptyUsers";
import { generateTableHeaders, combineDataSets } from "./UsersTableConfig";
import DeleteUserForm from "./components/DeleteUserForm";
import ResetPasswordModal from "./components/ResetPasswordModal";
import ResetSessionsModal from "./components/ResetSessionsModal";
import { NewUserType } from "./components/UserForm/UserForm";
import CreateUserModal from "./components/CreateUserModal";
import EditUserModal from "./components/EditUserModal";

const baseClass = "user-management";

interface IUserManagementProps {
  router: InjectedRouter; // v3
}

interface ITeamsResponse {
  teams: ITeam[];
}

const UserManagementPage = ({ router }: IUserManagementProps): JSX.Element => {
  const { config, currentUser, isPremiumTier } = useContext(AppContext);
  const { renderFlash } = useContext(NotificationContext);

  // STATES

  const [showCreateUserModal, setShowCreateUserModal] = useState<boolean>(
    false
  );
  const [showEditUserModal, setShowEditUserModal] = useState<boolean>(false);
  const [showDeleteUserModal, setShowDeleteUserModal] = useState<boolean>(
    false
  );
  const [showResetPasswordModal, setShowResetPasswordModal] = useState<boolean>(
    false
  );
  const [showResetSessionsModal, setShowResetSessionsModal] = useState<boolean>(
    false
  );
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isEditingUser, setIsEditingUser] = useState<boolean>(false);
  const [userEditing, setUserEditing] = useState<any>(null);
  const [createUserErrors, setCreateUserErrors] = useState<IUserFormErrors>(
    DEFAULT_CREATE_USER_ERRORS
  );
  const [editUserErrors, setEditUserErrors] = useState<IUserFormErrors>(
    DEFAULT_CREATE_USER_ERRORS
  );
  const [querySearchText, setQuerySearchText] = useState<string>("");

  // API CALLS
  const {
    data: teams,
    isFetching: isFetchingTeams,
    error: loadingTeamsError,
  } = useQuery<ITeamsResponse, Error, ITeam[]>(
    ["teams"],
    () => teamsAPI.loadAll(),
    {
      enabled: !!isPremiumTier,
      select: (data: ITeamsResponse) => data.teams,
    }
  );

  const {
    data: users,
    isFetching: isFetchingUsers,
    error: loadingUsersError,
    refetch: refetchUsers,
  } = useQuery<IUser[], Error, IUser[]>(
    ["users", querySearchText],
    () => usersAPI.loadAll({ globalFilter: querySearchText }),
    {
      select: (data: IUser[]) => data,
    }
  );

  const {
    data: invites,
    isFetching: isFetchingInvites,
    error: loadingInvitesError,
    refetch: refetchInvites,
  } = useQuery<IInvite[], Error, IInvite[]>(
    ["invites", querySearchText],
    () => invitesAPI.loadAll({ globalFilter: querySearchText }),
    {
      select: (data: IInvite[]) => {
        return data;
      },
    }
  );

  // TOGGLE MODALS

  const toggleCreateUserModal = useCallback(() => {
    setShowCreateUserModal(!showCreateUserModal);

    // clear errors on close
    if (!showCreateUserModal) {
      setCreateUserErrors(DEFAULT_CREATE_USER_ERRORS);
    }
  }, [showCreateUserModal, setShowCreateUserModal]);

  const toggleDeleteUserModal = useCallback(
    (user?: IUser | IInvite) => {
      setShowDeleteUserModal(!showDeleteUserModal);
      setUserEditing(!showDeleteUserModal ? user : null);
    },
    [showDeleteUserModal, setShowDeleteUserModal, setUserEditing]
  );

  const toggleEditUserModal = useCallback(
    (user?: IUser | IInvite) => {
      setShowEditUserModal(!showEditUserModal);
      setUserEditing(!showEditUserModal ? user : null);
      setEditUserErrors(DEFAULT_CREATE_USER_ERRORS);
    },
    [showEditUserModal, setShowEditUserModal, setUserEditing]
  );

  const toggleResetPasswordUserModal = useCallback(
    (user?: IUser | IInvite) => {
      setShowResetPasswordModal(!showResetPasswordModal);
      setUserEditing(!showResetPasswordModal ? user : null);
    },
    [showResetPasswordModal, setShowResetPasswordModal, setUserEditing]
  );

  const toggleResetSessionsUserModal = useCallback(
    (user?: IUser | IInvite) => {
      setShowResetSessionsModal(!showResetSessionsModal);
      setUserEditing(!showResetSessionsModal ? user : null);
    },
    [showResetSessionsModal, setShowResetSessionsModal, setUserEditing]
  );

  // FUNCTIONS

  const combineUsersAndInvites = memoize(
    (usersData, invitesData, currentUserId) => {
      return combineDataSets(usersData, invitesData, currentUserId);
    }
  );

  const goToUserSettingsPage = () => {
    const { USER_SETTINGS } = paths;
    router.push(USER_SETTINGS);
  };

  // NOTE: this is called once on the initial rendering. The initial render of
  // the TableContainer child component calls this handler.
  const onTableQueryChange = (queryData: ITableQueryData) => {
    const { searchQuery, sortHeader, sortDirection } = queryData;
    let sortBy: any = []; // TODO
    if (sortHeader !== "") {
      sortBy = [{ id: sortHeader, direction: sortDirection }];
    }

    setQuerySearchText(searchQuery);

    refetchUsers();
    refetchInvites();
  };

  const onActionSelect = (value: string, user: IUser | IInvite) => {
    switch (value) {
      case "edit":
        toggleEditUserModal(user);
        break;
      case "delete":
        toggleDeleteUserModal(user);
        break;
      case "passwordReset":
        toggleResetPasswordUserModal(user);
        break;
      case "resetSessions":
        toggleResetSessionsUserModal(user);
        break;
      case "editMyAccount":
        goToUserSettingsPage();
        break;
      default:
        return null;
    }
    return null;
  };

  const getUser = (type: string, id: number) => {
    let userData;
    if (type === "user") {
      userData = users?.find((user) => user.id === id);
    } else {
      userData = invites?.find((invite) => invite.id === id);
    }
    return userData;
  };

  const onCreateUserSubmit = (formData: any) => {
    setIsLoading(true);

    if (formData.newUserType === NewUserType.AdminInvited) {
      // Do some data formatting adding `invited_by` for the request to be correct and deleteing uncessary fields
      const requestData = {
        ...formData,
        invited_by: formData.currentUserId,
      };
      delete requestData.currentUserId; // this field is not needed for the request
      delete requestData.newUserType; // this field is not needed for the request
      delete requestData.password; // this field is not needed for the request
      invitesAPI
        .create(requestData)
        .then(() => {
          renderFlash(
            "success",
            `An invitation email was sent from ${config?.smtp_settings.sender_address} to ${formData.email}.`
          );
          toggleCreateUserModal();
          refetchInvites();
        })
        .catch((userErrors: { data: IApiError }) => {
          if (userErrors.data.errors[0].reason.includes("already exists")) {
            setCreateUserErrors({
              email: "A user with this email address already exists",
            });
          } else if (
            userErrors.data.errors[0].reason.includes("required criteria")
          ) {
            setCreateUserErrors({
              password: "Password must meet the criteria below",
            });
          } else {
            renderFlash("error", "Could not create user. Please try again.");
          }
        })
        .finally(() => {
          setIsLoading(false);
        });
    } else {
      // Do some data formatting deleting unnecessary fields
      const requestData = {
        ...formData,
      };
      delete requestData.currentUserId; // this field is not needed for the request
      delete requestData.newUserType; // this field is not needed for the request
      usersAPI
        .createUserWithoutInvitation(requestData)
        .then(() => {
          renderFlash("success", `Successfully created ${requestData.name}.`);
          toggleCreateUserModal();
          refetchUsers();
        })
        .catch((userErrors: { data: IApiError }) => {
          if (userErrors.data.errors[0].reason.includes("Duplicate")) {
            setCreateUserErrors({
              email: "A user with this email address already exists",
            });
          } else if (
            userErrors.data.errors[0].reason.includes("required criteria")
          ) {
            setCreateUserErrors({
              password: "Password must meet the criteria below",
            });
          } else {
            renderFlash("error", "Could not create user. Please try again.");
          }
        })
        .finally(() => {
          setIsLoading(false);
        });
    }
  };

  const onEditUser = (formData: any) => {
    const userData = getUser(userEditing.type, userEditing.id);

    setIsEditingUser(true);
    if (userEditing.type === "invite") {
      return (
        userData &&
        invitesAPI
          .update(userData.id, formData)
          .then(() => {
            renderFlash("success", `Successfully edited ${userEditing?.name}`);
            toggleEditUserModal();
            refetchInvites();
          })
          .catch((userErrors: { data: IApiError }) => {
            if (userErrors.data.errors[0].reason.includes("already exists")) {
              setEditUserErrors({
                email: "A user with this email address already exists",
              });
            } else if (
              userErrors.data.errors[0].reason.includes("required criteria")
            ) {
              setEditUserErrors({
                password: "Password must meet the criteria below",
              });
            } else {
              renderFlash(
                "error",
                `Could not edit ${userEditing?.name}. Please try again.`
              );
            }
          })
          .finally(() => {
            setIsEditingUser(false);
          })
      );
    }

    if (currentUser?.id === userEditing.id) {
      return (
        userData &&
        usersAPI
          .update(userData.id, formData)
          .then(() => {
            renderFlash("success", `Successfully edited ${userEditing?.name}`);
            toggleEditUserModal();
            refetchUsers();
          })
          .catch((userErrors: { data: IApiError }) => {
            if (userErrors.data.errors[0].reason.includes("already exists")) {
              setEditUserErrors({
                email: "A user with this email address already exists",
              });
            } else if (
              userErrors.data.errors[0].reason.includes("required criteria")
            ) {
              setEditUserErrors({
                password: "Password must meet the criteria below",
              });
            }
            renderFlash(
              "error",
              `Could not edit ${userEditing?.name}. Please try again.`
            );
          })
      );
    }

    let userUpdatedFlashMessage = `Successfully edited ${formData.name}`;

    if (userData?.email !== formData.email) {
      userUpdatedFlashMessage += `: A confirmation email was sent from ${config?.smtp_settings.sender_address} to ${formData.email}`;
    }

    return (
      userData &&
      usersAPI
        .update(userData.id, formData)
        .then(() => {
          renderFlash("success", userUpdatedFlashMessage);
          toggleEditUserModal();
          refetchUsers();
        })
        .catch((userErrors: { data: IApiError }) => {
          if (userErrors.data.errors[0].reason.includes("already exists")) {
            setEditUserErrors({
              email: "A user with this email address already exists",
            });
          } else if (
            userErrors.data.errors[0].reason.includes("required criteria")
          ) {
            setEditUserErrors({
              password: "Password must meet the criteria below",
            });
          } else {
            renderFlash(
              "error",
              `Could not edit ${userEditing?.name}. Please try again.`
            );
          }
        })
    );
  };

  const onDeleteUser = () => {
    if (userEditing.type === "invite") {
      invitesAPI
        .destroy(userEditing.id)
        .then(() => {
          renderFlash("success", `Successfully deleted ${userEditing?.name}.`);
        })
        .catch(() => {
          renderFlash(
            "error",
            `Could not delete ${userEditing?.name}. Please try again.`
          );
        })
        .finally(() => {
          toggleDeleteUserModal();
          refetchInvites();
        });
    } else {
      usersAPI
        .destroy(userEditing.id)
        .then(() => {
          renderFlash("success", `Successfully deleted ${userEditing?.name}.`);
        })
        .catch(() => {
          renderFlash(
            "error",
            `Could not delete ${userEditing?.name}. Please try again.`
          );
        })
        .finally(() => {
          toggleDeleteUserModal();
          refetchUsers();
        });
    }
  };

  const onResetSessions = () => {
    const isResettingCurrentUser = currentUser?.id === userEditing.id;

    usersAPI
      .deleteSessions(userEditing.id)
      .then(() => {
        if (isResettingCurrentUser) {
          clearToken();
          setTimeout(() => {
            window.location.href = "/";
          }, 500);
          return;
        }
        renderFlash("success", "Successfully reset sessions.");
      })
      .catch(() => {
        renderFlash("error", "Could not reset sessions. Please try again.");
      })
      .finally(() => {
        toggleResetSessionsUserModal();
      });
  };

  const resetPassword = (user: IUser) => {
    return usersAPI
      .requirePasswordReset(user.id, { require: true })
      .then(() => {
        renderFlash("success", "Successfully required a password reset.");
      })
      .catch(() => {
        renderFlash(
          "error",
          "Could not require a password reset. Please try again."
        );
      })
      .finally(() => {
        toggleResetPasswordUserModal();
      });
  };

  const renderEditUserModal = () => {
    const userData = getUser(userEditing.type, userEditing.id);

    return (
      <Modal
        title="Edit user"
        onExit={toggleEditUserModal}
        className={`${baseClass}__edit-user-modal`}
      >
        <>
          <EditUserModal
            defaultEmail={userData?.email}
            defaultName={userData?.name}
            defaultGlobalRole={userData?.global_role}
            defaultTeams={userData?.teams}
            onCancel={toggleEditUserModal}
            onSubmit={onEditUser}
            availableTeams={teams || []}
            isPremiumTier={isPremiumTier || false}
            smtpConfigured={config?.smtp_settings.configured || false}
            canUseSso={config?.sso_settings.enable_sso || false}
            isSsoEnabled={userData?.sso_enabled}
            isModifiedByGlobalAdmin
            isInvitePending={userEditing.type === "invite"}
            editUserErrors={editUserErrors}
            isLoading={isEditingUser}
          />
        </>
      </Modal>
    );
  };

  const renderCreateUserModal = () => {
    return (
      <CreateUserModal
        createUserErrors={createUserErrors}
        onCancel={toggleCreateUserModal}
        onSubmit={onCreateUserSubmit}
        availableTeams={teams || []}
        defaultGlobalRole={"observer"}
        defaultTeams={[]}
        isPremiumTier={isPremiumTier || false}
        smtpConfigured={config?.smtp_settings.configured || false}
        canUseSso={config?.sso_settings.enable_sso || false}
        isLoading={isLoading}
        isModifiedByGlobalAdmin
      />
    );
  };

  const renderDeleteUserModal = () => {
    return (
      <Modal
        title={"Delete user"}
        onExit={toggleDeleteUserModal}
        className={`${baseClass}__delete-user-modal`}
      >
        <DeleteUserForm
          name={userEditing.name}
          onDelete={onDeleteUser}
          onCancel={toggleDeleteUserModal}
        />
      </Modal>
    );
  };

  const renderResetPasswordModal = () => {
    return (
      <ResetPasswordModal
        user={userEditing}
        modalBaseClass={baseClass}
        onResetConfirm={resetPassword}
        onResetCancel={toggleResetPasswordUserModal}
      />
    );
  };

  const renderResetSessionsModal = () => {
    return (
      <ResetSessionsModal
        user={userEditing}
        modalBaseClass={baseClass}
        onResetConfirm={onResetSessions}
        onResetCancel={toggleResetSessionsUserModal}
      />
    );
  };

  const tableHeaders = generateTableHeaders(
    onActionSelect,
    isPremiumTier || false
  );

  const loadingTableData =
    isFetchingUsers || isFetchingInvites || isFetchingTeams;
  const tableDataError =
    loadingUsersError || loadingInvitesError || loadingTeamsError;

  let tableData: unknown = [];
  if (!loadingTableData) {
    tableData = combineUsersAndInvites(users, invites, currentUser?.id);
  }

  return (
    <div className={`${baseClass} body-wrap`}>
      <p className={`${baseClass}__page-description`}>
        Create new users, customize user permissions, and remove users from
        Fleet.
      </p>
      {/* TODO: find a way to move these controls into the table component */}
      {tableDataError ? (
        <TableDataError />
      ) : (
        <TableContainer
          columns={tableHeaders}
          data={tableData}
          isLoading={loadingTableData}
          defaultSortHeader={"name"}
          defaultSortDirection={"asc"}
          inputPlaceHolder={"Search"}
          actionButtonText={"Create user"}
          onActionButtonClick={toggleCreateUserModal}
          onQueryChange={onTableQueryChange}
          resultsTitle={"users"}
          emptyComponent={EmptyUsers}
          searchable
          showMarkAllPages={false}
          isAllPagesSelected={false}
          isClientSidePagination
        />
      )}
      {showCreateUserModal && renderCreateUserModal()}
      {showEditUserModal && renderEditUserModal()}
      {showDeleteUserModal && renderDeleteUserModal()}
      {showResetSessionsModal && renderResetSessionsModal()}
      {showResetPasswordModal && renderResetPasswordModal()}
    </div>
  );
};

export default UserManagementPage;
