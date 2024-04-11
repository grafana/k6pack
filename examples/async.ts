import { User, newUser } from "./user";

export default async () => {
  const user: User = newUser("John");
  console.log(user);
};
