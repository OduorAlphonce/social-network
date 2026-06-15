import { useCallback, useEffect, useMemo, useState } from "react";
import { apiFetch, ApiError } from "../utils/api";
import { AuthContext } from "./auth-context";

/**
 * Provides shared authentication state and actions to descendant components.
 *
 * @param {{children: import("react").ReactNode}} props Provider properties.
 * @returns {import("react").JSX.Element} Authentication context provider.
 */
export const AuthProvider = ({ children }) => {
  const [currentUser, setCurrentUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const refresh = useCallback(async () => {
    setIsLoading(true);

    try {
      const user = await apiFetch("/api/users/me");
      setCurrentUser(user);
      return user;
    } catch (error) {
      setCurrentUser(null);

      if (error instanceof ApiError && error.status === 401) {
        return null;
      }

      throw error;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const login = useCallback(
    async (credentials) => {
      setIsLoading(true);

      try {
        await apiFetch("/api/users/login", {
          method: "POST",
          body: credentials,
        });
        return await refresh();
      } catch (error) {
        setCurrentUser(null);
        setIsLoading(false);
        throw error;
      }
    },
    [refresh]
  );

  const logout = useCallback(async () => {
    setIsLoading(true);

    try {
      await apiFetch("/api/users/logout", { method: "POST" });
      setCurrentUser(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    let isActive = true;

    const loadCurrentUser = async () => {
      try {
        const user = await apiFetch("/api/users/me");

        if (isActive) {
          setCurrentUser(user);
        }
      } catch {
        if (isActive) {
          setCurrentUser(null);
        }
      } finally {
        if (isActive) {
          setIsLoading(false);
        }
      }
    };

    loadCurrentUser();

    return () => {
      isActive = false;
    };
  }, []);

  const value = useMemo(
    () => ({
      currentUser,
      isAuthenticated: currentUser !== null,
      isLoading,
      login,
      logout,
      refresh,
    }),
    [currentUser, isLoading, login, logout, refresh]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
