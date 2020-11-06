import { Box, Nav } from 'grommet';
import React from 'react';
import SidebarButton from '../SidebarButton/SidebarButton';
import { useHistory, useParams } from 'react-router-dom';

const Layout = ({ children }) => {
  const history = useHistory();
  const [active, setActive] = React.useState();
  const { patientId } = useParams();
  return (
    <Box
      align={'start'}
      justify={'start'}
      direction={'row-responsive'}
      gap={'medium'}
      pad={'medium'}
      fill={'horizontal'}
    >
      <Nav background={'neutral-2'}>
        {[
          { label: 'General Info', link: `/patient/${patientId}` },
          {
            label: 'Current Pregnancy',
            link: `/patient/${patientId}/current_pregnancy`,
          },
          { label: 'Lab Results', link: `/patient/${patientId}/lab_results` },
          { label: 'Home Visits', link: `/patient/${patientId}/home_visits` },
          {
            label: 'Arvs',
            link: `/patient/${patientId}/arvs`,
          },
          {
            label: 'Infant Diagnoses',
            link: `/patient/${patientId}/infant/diagnoses`,
          },
          {
            label: 'Hospital/Clinic Admissions',
            link: `/patient/${patientId}/admissions`,
          },
          {
            label: 'HIV Screenings',
            link: `/patient/${patientId}/hiv_screenings`,
          },
          {
            label: 'Contraceptives',
            link: `/patient/${patientId}/contraceptives`,
          },
        ].map((d) => {
          return (
            <SidebarButton
              key={d.label}
              label={d.label}
              active={d.label === active}
              link={d.link}
              onClick={() => {
                setActive(d.label);
                history.push(d.link, { id: patientId });
              }}
            />
          );
        })}
      </Nav>
      {children}
    </Box>
  );
};

export default Layout;
