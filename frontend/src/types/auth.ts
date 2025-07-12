
import { User } from "@/types/user";


// 或者使用更标准的命名方式
export interface AuthContextValue {
  user: User;
  setUser: (user: User) => void;
}