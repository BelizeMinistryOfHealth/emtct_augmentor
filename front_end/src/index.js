import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import * as serviceWorker from './serviceWorker';
import { Auth0Provider } from '@auth0/auth0-react';
import { RecoilRoot } from 'recoil';

const auth0Domain = 'emtct-dev.us.auth0.com';
const auth0ClientId = 'k46hfbBUDsOaPgNU9IlUd7hoWJ5Ku0EB';
const redirectUrl = window.location.origin;

ReactDOM.render(
  <React.StrictMode>
    <Auth0Provider
      domain={auth0Domain}
      clientId={auth0ClientId}
      redirectUri={redirectUrl}
    >
      <RecoilRoot>
        <App />
      </RecoilRoot>
    </Auth0Provider>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
