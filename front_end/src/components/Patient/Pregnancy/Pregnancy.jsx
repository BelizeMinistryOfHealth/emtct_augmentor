import React from 'react';
import Layout from '../../Layout/Layout';
import PregnancyVitals from './PregnancyVitals/PregnancyVitals';

const Pregnancy = (props) => {
  console.dir({ pregProps: props });
  const { location } = props;
  return (
    <Layout location={location} props={props}>
      <PregnancyVitals />
    </Layout>
  );
};

export default Pregnancy;
