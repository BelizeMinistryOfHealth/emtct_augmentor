import React from 'react';
import { Route, Switch } from 'react-router-dom';
import LabResults from './components/Patient/LabResults/LabResults';
import ArvTreatment from './components/Patient/ArvTreatment/ArvTreatment';
import SyphilisTreatment from './components/Patient/SyphilisTreatment/SyphilisTreatment';
import InfantDiagnoses from './components/Infant/Diagnoses';
import InfantSyphillisTreatment from './components/Infant/SyphillisTreatment';
import InfantSyphilisScreenings from './components/Infant/SyphilisScreenings';
import HivScreeningCreateForm from './components/Infant/HivScreenings/HivScreeningCreate';
import InfantHivScreenings from './components/Infant/HivScreenings';
import Infant from './components/Infant';
import HomeVisitCreateForm from './components/Patient/HomeVisit/HomeVisitCreate';
import HomeVisitList from './components/Patient/HomeVisit';
import HospitalAdmissionCreateForm from './components/Patient/HospitalAdmissions/HospitalAdmissionsCreate';
import HospitalAdmissions from './components/Patient/HospitalAdmissions';
import ContraceptivesCreateForm from './components/Patient/Contraceptives/ContraceptivesCreate';
import ContraceptivesUsed from './components/Patient/Contraceptives';
import PartnerSyphilisTreatmentCreate from './components/Partner/SyphilisTreatments/PartnerSyphilisTreatmentCreate';
import PartnerSyphilisTreatments from './components/Partner/SyphilisTreatments/PartnerSyphilisTreatments';
import CcontactTracingCreate from './components/Partner/ContactTracing/ContactTracingCreate';
import ContactTracing from './components/Partner/ContactTracing/ContactTracing';
import CurrentPregnancy from './components/Patient/Pregnancy/CurrentPregnancy/CurrentPregnancy';
import Overview from './components/Patient/Overview';
import ReportHome from './components/Reports/Home';
import UserList from './components/Users/UserList';
import Search from './components/Search/Search';
import { Main } from 'grommet';
import InfantPcrs from './components/Reports/InfantPcrs/InfantPcrs';

const Router = () => {
  return (
    <Main>
      <Switch>
        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/lab_results'}
          component={LabResults}
        />

        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/arvs'}
          component={ArvTreatment}
        />
        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/syphilisTreatment'}
          component={SyphilisTreatment}
        />
        <Route
          path={
            '/patient/:patientId/pregnancy/:pregnancyId/infant/:infantId/diagnoses'
          }
          component={InfantDiagnoses}
        />
        <Route
          path={
            '/patient/:patientId/pregnancy/:pregnancyId/infant/:infantId/syphilisTreatment'
          }
          component={InfantSyphillisTreatment}
        />
        <Route
          path={
            '/patient/:patientId/pregnancy/:pregnancyId/infant/:infantId/syphilisScreenings'
          }
          component={InfantSyphilisScreenings}
        />
        <Route
          path={
            '/patient/:patientId/pregnancy/:pregnancyId/infant/:infantId/hivScreenings/new'
          }
          component={HivScreeningCreateForm}
        />
        <Route
          path={
            '/patient/:patientId/pregnancy/:pregnancyId/infant/:infantId/hivScreenings'
          }
          component={InfantHivScreenings}
        />

        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/infant'}
          component={Infant}
        />
        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/home_visits/new'}
          component={HomeVisitCreateForm}
        />
        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/home_visits'}
          component={HomeVisitList}
        />
        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/admissions/new'}
          component={HospitalAdmissionCreateForm}
        />
        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/admissions'}
          component={HospitalAdmissions}
        />
        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/contraceptives/new'}
          component={ContraceptivesCreateForm}
        />
        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId/contraceptives'}
          component={ContraceptivesUsed}
        />
        <Route
          path={
            '/patient/:patientId/pregnancy/:pregnancyId/partners/syphilisTreatments/new'
          }
          component={PartnerSyphilisTreatmentCreate}
        />
        <Route
          path={
            '/patient/:patientId/pregnancy/:pregnancyId/partners/syphilisTreatments'
          }
          component={PartnerSyphilisTreatments}
        />
        <Route
          path={
            '/patient/:patientId/pregnancy/:pregnancyId/partners/contactTracing/new'
          }
          component={CcontactTracingCreate}
        />
        <Route
          path={
            '/patient/:patientId/pregnancy/:pregnancyId/partners/contactTracing'
          }
          component={ContactTracing}
        />

        <Route
          path={'/patient/:patientId/pregnancy/:pregnancyId'}
          component={CurrentPregnancy}
        />

        <Route exact={true} path={'/patient/:patientId'} component={Overview} />
        <Route path={'/reports/pcrs'} component={InfantPcrs} />
        <Route path={'/reports'} component={ReportHome} />
        <Route exact={true} path={'/admin/users'} component={UserList} />
        <Route exact={true} path={'/search'} component={Search} />
        <Route exact={true} path={'/'} component={ReportHome} />
      </Switch>
    </Main>
  );
};

export default Router;
