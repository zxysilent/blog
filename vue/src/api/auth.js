import ajax from "./base";

export const login = data => {
  return ajax.post("/login", data);
};
export const auth = () => {
  return ajax.get("/auth");
};
export const logoff = data => {
  //   return ajax.post("/auth/login", data);
};
//统计
export const stat = () => {
  return ajax.get("/stat");
};
