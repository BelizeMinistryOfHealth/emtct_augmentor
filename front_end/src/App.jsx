import React from 'react';
import { Box, Grommet, Main } from 'grommet';
import { grommet } from 'grommet/themes';
import Navbar from './components/Navbar/Navbar';
import Search from './components/Search/Search';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import CurrentPregnancy from './components/Patient/Pregnancy/CurrentPregnancy/CurrentPregnancy';
import HttpApiProvider from './providers/HttpProvider';
import HomeVisitList from './components/Patient/HomeVisit';
import HomeVisitCreateForm from './components/Patient/HomeVisit/HomeVisitCreate';
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
import HivScreeningCreateForm from './components/Infant/HivScreenings/HivScreeningCreate';
import InfantSyphillisTreatment from './components/Infant/SyphillisTreatment';
import InfantSyphilisScreenings from './components/Infant/SyphilisScreenings';
import PartnerSyphilisTreatments from './components/Partner/SyphilisTreatments/PartnerSyphilisTreatments';
import PartnerSyphilisTreatmentCreate from './components/Partner/SyphilisTreatments/PartnerSyphilisTreatmentCreate';
import ContactTracing from './components/Partner/ContactTracing/ContactTracing';
import CcontactTracingCreate from './components/Partner/ContactTracing/ContactTracingCreate';
import { FirebaseAuthConsumer } from '@react-firebase/auth';
import Login from './components/Auth/Login';
import Overview from './components/Patient/Overview';

const { REACT_APP_API_URL } = process.env;

function App() {
  const [idToken, setIdToken] = React.useState();
  return (
    <FirebaseAuthConsumer>
      {(obj) => {
        const { isSignedIn, user } = obj;
        if (isSignedIn) {
          user.getIdToken().then(setIdToken);
          if (!idToken) {
            return <>Loading</>;
          }
          return (
            <Grommet theme={grommet}>
              <BrowserRouter>
                <Navbar />
                <HttpApiProvider idToken={idToken} baseUrl={REACT_APP_API_URL}>
                  <Main>
                    <Switch>
                      <Route
                        path={'/patient/:patientId/pregnancy/:pregnancyId'}
                        component={CurrentPregnancy}
                      />
                      <Route
                        path={'/patient/:patientId/home_visits/new'}
                        component={HomeVisitCreateForm}
                      />
                      <Route
                        path={'/patient/:patientId/home_visits'}
                        component={HomeVisitList}
                      />
                      <Route
                        path={'/patient/:patientId/contraceptives/new'}
                        component={ContraceptivesCreateForm}
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
                        path={
                          '/patient/:patientId/infant/:infantId/hivScreenings/new'
                        }
                        component={HivScreeningCreateForm}
                      />
                      <Route
                        path={
                          '/patient/:patientId/infant/:infantId/hivScreenings'
                        }
                        component={InfantHivScreenings}
                      />

                      <Route
                        path={'/patient/:patientId/infant/:infantId/diagnoses'}
                        component={InfantDiagnoses}
                      />
                      <Route
                        path={
                          '/patient/:patientId/infant/:infantId/syphilisTreatment'
                        }
                        component={InfantSyphillisTreatment}
                      />
                      <Route
                        path={
                          '/patient/:patientId/infant/:infantId/syphilisScreenings'
                        }
                        component={InfantSyphilisScreenings}
                      />
                      <Route
                        path={
                          '/patient/:patientId/partners/syphilisTreatments/new'
                        }
                        component={PartnerSyphilisTreatmentCreate}
                      />
                      <Route
                        path={'/patient/:patientId/partners/syphilisTreatments'}
                        component={PartnerSyphilisTreatments}
                      />
                      <Route
                        path={'/patient/:patientId/partners/contactTracing/new'}
                        component={CcontactTracingCreate}
                      />
                      <Route
                        path={'/patient/:patientId/partners/contactTracing'}
                        component={ContactTracing}
                      />
                      <Route
                        path={'/patient/:patientId/infant'}
                        component={Infant}
                      />
                      <Route
                        path={'/patient/:patientId'}
                        component={Overview}
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
