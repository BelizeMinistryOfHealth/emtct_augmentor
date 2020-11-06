import { Card } from 'grommet';
import React from 'react';

const AppCard = ({ children, ...rest }) => {
  return (
    <Card background={'neutral-3'} {...rest}>
      {children}
    </Card>
  );
};

export default AppCard;
