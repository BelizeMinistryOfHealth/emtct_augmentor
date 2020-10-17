import React from 'react';
import { Grommet, Main } from 'grommet';
import { grommet } from 'grommet/themes';
import { useAuth0 } from '@auth0/auth0-react';
import Navbar from './components/Navbar/Navbar';
import Welcome from './components/Welcome/Welcome';
import Search from './components/Search/Search';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import PatientSummary from './components/Patient/PatientSummary/PatientSummary.jsx';
import CurrentPregnancy from './components/Patient/Pregnancy/CurrentPregnancy/CurrentPregnancy';

function App() {
  const { isAuthenticated, getIdTokenClaims } = useAuth0();

  // TODO: Remove this. It is here as a hack for getting a valid access token when we need to test the backend.
  React.useEffect(() => {
    if (isAuthenticated) {
      console.log('getting token');
      (async () => {
        const idToken = await getIdTokenClaims({
          audience: 'https://emtct-dev.us.auth0.com/userinfo',
          scope: 'read:hiv',
        });
        console.dir({ idToken });
      })();
    }
    console.log('not authenticated');
  }, [isAuthenticated, getIdTokenClaims]);

  const fullTheme = !isAuthenticated;

  if (isAuthenticated) {
    return (
      <Grommet theme={grommet} full={fullTheme}>
        <BrowserRouter>
          <Navbar />

          <Main>
            <Switch>
              <Route
                path={'/patient/:id/current_pregnancy'}
                component={CurrentPregnancy}
              />
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
