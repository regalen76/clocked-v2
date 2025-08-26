import * as z from "zod";

export const loginSchema = z.object({
  email: z.email(),
  password: z.string().min(8),
});

export type Login = z.infer<typeof loginSchema>;

export interface LoginResponse {
  data: string;
  message: string;
  status: string;
}

export interface User {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: null | string;
  username: string;
  email: string;
  names: string;
  avatar: string;
}

export interface UserResponse {
  data: User;
  message: string;
  status: string;
}

