import { Box, Nav } from 'grommet';
import React from 'react';
import SidebarButton from '../SidebarButton/SidebarButton';

const Layout = ({ children, ...rest }) => {
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
          'General Info',
          'Current Pregnancy',
          'Lab Results',
          'Home Visits',
          'Hospital/Clinic Admissions',
        ].map((label) => {
          return (
            <SidebarButton
              key={label}
              label={label}
              active={label === active}
              onClick={() => setActive(label)}
            />
          );
        })}
      </Nav>
      {children}
    </Box>
  );
};

export default Layout;
