import React, { useEffect } from 'react';
import { Grommet, Main, Nav } from 'grommet';
import { grommet } from 'grommet/themes';
import { useAuth0, withAuth0 } from '@auth0/auth0-react';
import Navbar from './components/Navbar/Navbar';
import Welcome from './components/Welcome/Welcome';
import Search from './components/Search/Search';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import PatientSummary from './components/Patient/PatientSummary/PatientSummary.jsx';

const Profile = (props) => {
  const {
    user,
    isAuthenticated,
    getAccessTokenSilently,
    getIdTokenClaims,
  } = props.auth0;

  useEffect(() => {
    if (isAuthenticated) {
      console.log('getting token');
      (async () => {
        const token = await getAccessTokenSilently({
          audience: 'https://emtct-dev.us.auth0.com/userinfo',
          scope: 'read:hiv',
        });
        console.log({ token });
        const idToken = await getIdTokenClaims({
          audience: 'https://emtct-dev.us.auth0.com/userinfo',
          scope: 'read:hiv',
        });
        console.dir({ idToken });
      })();
    }
    console.log('not authenticated');
  }, [getAccessTokenSilently, isAuthenticated, getIdTokenClaims]);

  return isAuthenticated && <div>Hello {user.name}</div>;
};

const ProfileComponent = withAuth0(Profile);

function App() {
  const { isAuthenticated } = useAuth0();

  const fullTheme = !isAuthenticated;

  if (isAuthenticated) {
    return (
      <Grommet theme={grommet} full={fullTheme}>
        <BrowserRouter>
          <Navbar />

          <Main>
            <Switch>
              <Route path={'/patient/:id'} component={PatientSummary} />
              <Route path={'/'} component={Search} />
            </Switch>
          </Main>
        </BrowserRouter>
      </Grommet>
    );
  }
  return (
    <Grommet theme={grommet} full={fullTheme}>
      <Navbar />
      <Main>
        <Welcome />
      </Main>
    </Grommet>
  );
}

export default App;
