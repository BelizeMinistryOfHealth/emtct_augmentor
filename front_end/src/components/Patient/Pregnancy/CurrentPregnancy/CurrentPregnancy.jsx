import { Box } from 'grommet';
import React from 'react';
import { useRecoilValueLoadable } from 'recoil';
import { currentPregnancySelector } from '../../../../state';
import Layout from '../../../Layout/Layout';
import ArvTreatment from '../../ArvTreatment/ArvTreatment';
import PatientBasicInfo from '../../PatientBasicInfo/PatientBasicInfo';
import PregnancyVitals from '../PregnancyVitals/PregnancyVitals';
import PreNatalCare from '../PreNatalCare/PreNatalCare';

const CurrentPregnancy = (props) => {
  const { location } = props;

  const id = location.state.id;

  const { state, contents } = useRecoilValueLoadable(
    currentPregnancySelector(id)
  );
  let currentPregnancy = {};
  switch (state) {
    case 'hasValue':
      currentPregnancy = contents;
      break;
    case 'hasValue':
      return contents.message;
    case 'loading':
      return 'loading';
    default:
      return '';
  }

  return (
    <Layout location={location} props={props}>
      <Box
        direction={'column'}
        gap={'medium'}
        pad={'medium'}
        justify={'start'}
        align={'start'}
        fill
      >
        <Box
          direction={'row-responsive'}
          gap={'medium'}
          pad={'medium'}
          justify={'start'}
          align={'start'}
          fill
        >
          <PatientBasicInfo
            basicInfo={currentPregnancy.basicInfo}
            nextOfKin={currentPregnancy.nextOfKin}
          />
          <PregnancyVitals vitals={currentPregnancy.vitals} />
          <PreNatalCare info={currentPregnancy.prenatalCareInfo} />
        </Box>
        <Box
          gap={'medium'}
          pad={'medium'}
          justify={'center'}
          align={'center'}
          fill
        >
          <ArvTreatment
            patientId={currentPregnancy.basicInfo.patientId}
            encounterId={currentPregnancy.encounterId}
          />
        </Box>
      </Box>
    </Layout>
  );
};

export default CurrentPregnancy;
