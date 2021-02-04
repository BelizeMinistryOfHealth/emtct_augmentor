import React from 'react';
import { Box, Grommet, Main } from 'grommet';
import { grommet } from 'grommet/themes';
import Navbar from './components/Navbar/Navbar';
import { BrowserRouter } from 'react-router-dom';
import HttpApiProvider from './providers/HttpProvider';
import { FirebaseAuthConsumer } from '@react-firebase/auth';
import Login from './components/Auth/Login';
import jwt_decode from 'jwt-decode';
import Router from './routes';

const { REACT_APP_API_URL } = process.env;

function App() {
  const [token, setToken] = React.useState({
    idToken: null,
    permissions: [],
    exp: 0,
  });

  return (
    <FirebaseAuthConsumer>
      {(obj) => {
        const { isSignedIn, user } = obj;
        if (isSignedIn) {
          user.getIdToken().then((_idtoken) => {
            const decodedToken = jwt_decode(_idtoken);
            const newTokenDate = new Date(0);
            newTokenDate.setUTCSeconds(decodedToken.exp);
            if (!token.idToken || token.exp === 0) {
              const permissions = decodedToken.permissions;
              setToken({
                idToken: _idtoken,
                permissions,
                exp: decodedToken.exp,
              });
            }
          });

          if (!token.idToken) {
            return <>Loading</>;
          }
          return (
            <Grommet theme={grommet}>
              <HttpApiProvider
                idToken={token.idToken}
                baseUrl={REACT_APP_API_URL}
              >
                <BrowserRouter forceRefresh={false}>
                  <Navbar permissions={token.permissions} />
                  <Router />
                </BrowserRouter>
              </HttpApiProvider>
            </Grommet>
          );
        }

        return (
          <Grommet theme={grommet} full>
            <Main>
              <Box
                flex={'grow'}
                align={'center'}
                pad={'xxlarge'}
                background={'neutral-2'}
                fill
              >
                <Login />
              </Box>
            </Main>
          </Grommet>
        );
      }}
    </FirebaseAuthConsumer>
  );
}

export default App;
