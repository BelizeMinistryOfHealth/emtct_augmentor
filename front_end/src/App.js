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
import ContraceptivesUsed from './components/Patient/Contraceptives';
import ContraceptivesCreateForm from './components/Patient/Contraceptives/ContraceptivesCreate';
import HospitalAdmissions from './components/Patient/HospitalAdmissions';
import HospitalAdmissionCreateForm from './components/Patient/HospitalAdmissions/HospitalAdmissionsCreate';
import LabResults from './components/Patient/LabResults/LabResults';
import ArvTreatment from './components/Patient/ArvTreatment/ArvTreatment';
import SyphilisTreatment from './components/Patient/SyphilisTreatment/SyphilisTreatment';
import Infant from './components/Infant';
import InfantHivScreenings from './components/Infant/HivScreenings';
import InfantDiagnoses from './components/Infant/Diagnoses';

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
                  path={'/patient/:patientId/current_pregnancy'}
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
                  path={'/patient/:patientId/contraceptives/new'}
                  component={ContraceptivesCreateForm}
                />
                <Route
                  path={'/patient/:patientId/home_visits'}
                  component={HomeVisitList}
                />
                <Route
                  path={'/patient/:patientId/hiv_screenings'}
                  component={HivScreening}
                />
                <Route
                  path={'/patient/:patientId/contraceptives'}
                  component={ContraceptivesUsed}
                />
                <Route
                  path={'/patient/:patientId/admissions/new'}
                  component={HospitalAdmissionCreateForm}
                />
                <Route
                  path={'/patient/:patientId/admissions'}
                  component={HospitalAdmissions}
                />

                <Route
                  path={'/patient/:patientId/lab_results'}
                  component={LabResults}
                />
                <Route
                  path={'/patient/:patientId/arvs'}
                  component={ArvTreatment}
                />
                <Route
                  path={'/patient/:patientId/syphilisTreatment'}
                  component={SyphilisTreatment}
                />
                <Route
                  path={'/patient/:patientId/infant/:infantId/hivScreenings'}
                  component={InfantHivScreenings}
                />
                <Route
                  path={'/patient/:patientId/infant/:infantId/diagnoses'}
                  component={InfantDiagnoses}
                />
                <Route path={'/patient/:patientId/infant'} component={Infant} />
                <Route
                  path={'/patient/:patientId'}
                  component={PatientSummary}
                />
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
