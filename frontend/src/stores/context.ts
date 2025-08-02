import { createContext, useContext } from "react";
import { AuthContextValue } from "@/types/auth";
import { DefaultUser } from "@/types/user.ts";
import { initialState, ThemeProviderState } from "@/types/system.ts";

// 创建 AuthContext，初始值为 null
export const AuthContext = createContext<AuthContextValue>({ user: DefaultUser, setUser: () => {} });

// 创建自定义 hook 来使用 AuthContext
export const useAuth = (): AuthContextValue => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
};

export const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined) throw new Error("useTheme must be used within a ThemeProvider");

  return context;
};
