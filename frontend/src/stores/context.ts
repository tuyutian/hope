import {Context, createContext} from "react";
import {DefaultUser, User} from "@/types/user.ts";

export const AuthContext:Context<User> = createContext<User>(DefaultUser);

