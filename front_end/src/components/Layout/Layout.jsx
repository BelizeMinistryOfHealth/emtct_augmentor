import { Box, Nav } from 'grommet';
import React from 'react';
import SidebarButton from '../SidebarButton/SidebarButton';
import { useHistory } from 'react-router-dom';

const Layout = ({ children, location, ...rest }) => {
  const patientId = location?.state?.id;
  const history = useHistory();
  const [active, setActive] = React.useState();
  return (
    <Box
      align={'start'}
      justify={'start'}
      direction={'row-responsive'}
      gap={'medium'}
      pad={'medium'}
      fill
    >
      <Nav background={'neutral-2'}>
        {[
          { label: 'General Info', link: `/patient/${patientId}` },
          {
            label: 'Current Pregnancy',
            link: `/patient/${patientId}/current_pregnancy`,
          },
          { label: 'Lab Results', link: '/' },
          { label: 'Home Visits', link: `/patient/${patientId}/home_visits` },
          { label: 'Hospital/Clinic Admissions', link: '/' },
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
