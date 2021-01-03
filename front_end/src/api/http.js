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
        'Content-Type': 'application/json',
      },
    });
  }

  async get(path) {
    return this.axiosInstance.get(path);
  }

  async put(path, body) {
    return this.axiosInstance.put(path, JSON.stringify(body));
  }

  async post(path, body) {
    return this.axiosInstance.post(path, JSON.stringify(body));
  }

  async delete(path) {
    return this.axiosInstance.delete(path);
  }
}

export const useHttpApi = (accessToken, baseUrl) => {
  const [instance, setInstance] = useState();
  const i = new HttpApi(accessToken, baseUrl);
  setInstance(i);
  return instance;
};

export default HttpApi;
