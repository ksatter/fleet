import React, { useContext, useEffect, useState } from "react";
import { AxiosResponse } from "axios";
import { InjectedRouter } from "react-router";
import { QueryClient, QueryClientProvider } from "react-query";
import classnames from "classnames";

import TableProvider from "context/table";
import QueryProvider from "context/query";
import PolicyProvider from "context/policy";
import NotificationProvider from "context/notification";
import { AppContext } from "context/app";
import local, { authToken } from "utilities/local";
import useDeepEffect from "hooks/useDeepEffect";

import usersAPI from "services/entities/users";
import configAPI from "services/entities/config";

import { ErrorBoundary } from "react-error-boundary";
// @ts-ignore
import Fleet403 from "pages/errors/Fleet403";
// @ts-ignore
import Fleet404 from "pages/errors/Fleet404";
// @ts-ignore
import Fleet500 from "pages/errors/Fleet500";
import Spinner from "components/Spinner";

interface IAppProps {
  children: JSX.Element;
  router: InjectedRouter;
  location:
    | {
        pathname: string;
      }
    | undefined;
}

const App = ({ children, location, router }: IAppProps): JSX.Element => {
  const queryClient = new QueryClient();
  const {
    currentUser,
    isGlobalObserver,
    isOnlyObserver,
    isAnyTeamMaintainerOrTeamAdmin,
    setAvailableTeams,
    setCurrentUser,
    setConfig,
    setEnrollSecret,
  } = useContext(AppContext);

  const [isLoading, setIsLoading] = useState<boolean>(false);

  const fetchConfig = async () => {
    try {
      const config = await configAPI.loadAll();
      setConfig(config);
    } catch (error) {
      console.error(error);
      return false;
    } finally {
      setIsLoading(false);
    }
    return true;
  };

  const fetchCurrentUser = async () => {
    try {
      const { user, available_teams } = await usersAPI.me();
      setCurrentUser(user);
      setAvailableTeams(available_teams);
      fetchConfig();
    } catch (error) {
      if (!location?.pathname.includes("/login/reset")) {
        console.log(error);
        local.removeItem("auth_token");
        window.location.href = "/login";
      }
    }
    return true;
  };

  useEffect(() => {
    if (authToken()) {
      fetchCurrentUser();
    }
  }, [location?.pathname]);

  useDeepEffect(() => {
    const canGetEnrollSecret =
      currentUser &&
      typeof isGlobalObserver !== "undefined" &&
      !isGlobalObserver &&
      typeof isOnlyObserver !== "undefined" &&
      !isOnlyObserver &&
      typeof isAnyTeamMaintainerOrTeamAdmin !== "undefined" &&
      !isAnyTeamMaintainerOrTeamAdmin;

    const getEnrollSecret = async () => {
      try {
        const { spec } = await configAPI.loadEnrollSecret();
        setEnrollSecret(spec.secrets);
      } catch (error) {
        console.error(error);
        return false;
      }
    };

    if (canGetEnrollSecret) {
      getEnrollSecret();
    }
  }, [currentUser, isGlobalObserver, isOnlyObserver]);

  // "any" is used on purpose. We are using Axios but this
  // function expects a native React Error type, which is incompatible.
  const renderErrorOverlay = ({ error }: any) => {
    // @ts-ignore
    console.error(error);

    const overlayError = error as AxiosResponse;
    if (overlayError.status === 403 || overlayError.status === 402) {
      return <Fleet403 />;
    }

    if (overlayError.status === 404) {
      return <Fleet404 />;
    }

    return <Fleet500 />;
  };

  const wrapperStyles = classnames("wrapper");
  return isLoading ? (
    <Spinner />
  ) : (
    <QueryClientProvider client={queryClient}>
      <TableProvider>
        <QueryProvider>
          <PolicyProvider>
            <NotificationProvider>
              <ErrorBoundary
                fallbackRender={renderErrorOverlay}
                resetKeys={[location?.pathname]}
              >
                <div className={wrapperStyles}>{children}</div>
              </ErrorBoundary>
            </NotificationProvider>
          </PolicyProvider>
        </QueryProvider>
      </TableProvider>
    </QueryClientProvider>
  );
};

export default App;
