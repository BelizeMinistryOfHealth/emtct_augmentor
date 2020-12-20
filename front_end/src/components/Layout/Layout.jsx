import { Box, Nav } from 'grommet';
import React from 'react';
import SidebarButton from '../SidebarButton/SidebarButton';
import { useHistory, useParams } from 'react-router-dom';

const Layout = ({ children }) => {
  const history = useHistory();
  const [active, setActive] = React.useState();
  const { patientId, pregnancyId } = useParams();
  return (
    <Box
      align={'start'}
      justify={'start'}
      direction={'row-responsive'}
      gap={'xxsmall'}
      pad={'xxsmall'}
      fill={'horizontal'}
    >
      <Nav background={'neutral-2'}>
        {[
          {
            label: 'General Info',
            link: `/patient/${patientId}/pregnancy/${pregnancyId}`,
          },
          {
            label: 'Infant',
            link: `/patient/${patientId}/infant`,
          },
          {
            label: 'Partners',
            link: `/patient/${patientId}/partners/syphilisTreatments`,
          },
          {
            label: 'Lab Results',
            link: `/patient/${patientId}/pregnancy/${pregnancyId}/lab_results`,
          },
          { label: 'Home Visits', link: `/patient/${patientId}/home_visits` },
          {
            label: 'Arvs',
            link: `/patient/${patientId}/arvs`,
          },
          {
            label: 'Syphilis Treatment',
            link: `/patient/${patientId}/syphilisTreatment`,
          },
          {
            label: 'Hospital/Clinic Admissions',
            link: `/patient/${patientId}/admissions`,
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
