import React, { useRef } from 'react';
import HttpApi from '../api/http';

export const HttpApiContext = React.createContext();

export const useHttpApi = () => React.useContext(HttpApiContext);

const HttpApiProvider = ({ idToken, baseUrl, children }) => {
  const httpInstance = useRef();

  httpInstance.current = new HttpApi(idToken, baseUrl);

  const contextValue = { httpInstance: httpInstance.current };
  if (httpInstance.current) {
    return (
      <HttpApiContext.Provider value={contextValue}>
        {children}
      </HttpApiContext.Provider>
    );
  } else {
    return <>Loading...</>;
  }
};

export default HttpApiProvider;
