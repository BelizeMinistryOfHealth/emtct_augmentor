import { useState } from 'react';
import Axios from 'axios';

class HttpApi {
  accessToken;
  baseUrl;
  axiosInstance;
  constructor(accessToken, baseUrl) {
    this.accessToken = accessToken;
    this.baseUrl = baseUrl;
    this.axiosInstance = Axios.create({
      baseURL: this.baseUrl,
      headers: {
        Authorization: `Bearer ${accessToken}`,
        Accept: 'application/json',
      },
    });
  }

  async get(path) {
    return this.axiosInstance.get(path);
  }
}

export const useHttpApi = (accessToken, baseUrl) => {
  const [instance, setInstance] = useState();
  const i = new HttpApi(accessToken, baseUrl);
  setInstance(i);
  return instance;
};

export default HttpApi;
