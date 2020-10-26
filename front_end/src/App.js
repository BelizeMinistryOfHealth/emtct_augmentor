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
import HttpApiProvider from './providers/HttpProvider';
import HomeVisitList from './components/Patient/HomeVisit';
import HomeVisitCreateForm from './components/Patient/HomeVisit/HomeVisitCreate';
import HivScreening from './components/Patient/HivScreening';
import HivScreeningCreateForm from './components/Patient/HivScreening/HivScreeningCreate';

function App() {
  const { isAuthenticated, getIdTokenClaims } = useAuth0();
  const [idToken, setIdToken] = React.useState();
  const { REACT_APP_API_URL } = process.env;

  React.useEffect(() => {
    if (isAuthenticated) {
      console.log('getting token');
      (async () => {
        try {
          const idToken = await getIdTokenClaims();
          setIdToken(idToken.__raw);
        } catch (e) {
          console.error('error fetching token: ', e);
        }
      })();
    } else {
      console.log('not authenticated');
    }
  }, [isAuthenticated, getIdTokenClaims, setIdToken]);

  const fullTheme = !isAuthenticated;

  if (isAuthenticated && idToken) {
    return (
      <Grommet theme={grommet} full={fullTheme}>
        <BrowserRouter>
          <Navbar />
          <HttpApiProvider idToken={idToken} baseUrl={REACT_APP_API_URL}>
            <Main>
              <Switch>
                <Route
                  path={'/patient/:id/current_pregnancy'}
                  component={CurrentPregnancy}
                />
                <Route
                  path={'/patient/:patientId/home_visits/new'}
                  component={HomeVisitCreateForm}
                />
                <Route
                  path={'/patient/:patientId/hiv_screenings/new'}
                  component={HivScreeningCreateForm}
                />
                <Route
                  path={'/patient/:patientId/home_visits'}
                  component={HomeVisitList}
                />
                <Route
                  path={'/patient/:patientId/hiv_screenings'}
                  component={HivScreening}
                />
                <Route path={'/patient/:id'} component={PatientSummary} />
                <Route path={'/'} component={Search} />
              </Switch>
            </Main>
          </HttpApiProvider>
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
