import React from 'react';
import { Box, Card, CardBody, CardHeader, Text } from 'grommet';
import Search from '../Search';
import { useHttpApi } from '../../../providers/HttpProvider';
import InfantPcrsList from './InfantPcrsList';
import { mergeScreenings } from './missingPcrs';

const InfantPcrs = () => {
  const { httpInstance } = useHttpApi();
  const [year, setYear] = React.useState();
  const [pcrs, setPcrs] = React.useState({
    loading: false,
    data: [],
    error: undefined,
  });

  React.useEffect(() => {
    const search = async (year) => {
      try {
        const result = await httpInstance.get(`/reports/missingPcrs/${year}`);
        const screenings = mergeScreenings(result.data);
        setPcrs({ loading: false, data: screenings, error: undefined });
      } catch (e) {
        console.error(e);
        setPcrs({
          loading: false,
          data: [],
          error: 'Error retrieving missing PCRS',
        });
      }
    };
    if (pcrs.loading) {
      search(year);
    }
  }, [httpInstance, year, pcrs]);

  const onSubmit = (yr) => {
    setYear(yr);
    setPcrs({ ...pcrs, loading: true });
  };
  return (
    <Box direction={'column'} gap={'small'} pad={'small'}>
      <Card
        margin={{ top: 'medium', left: 'large', right: 'large' }}
        pad={'small'}
        gap={'xxsmall'}
        responsive={true}
      >
        <CardHeader justify={'center'}>
          <Text size={'large'} weight={'bold'}>
            Infant PCRs
          </Text>
        </CardHeader>
        <CardBody>
          <Search
            label={'Enter a Year'}
            onSubmit={onSubmit}
            errMessage={pcrs.error ?? ''}
          />
        </CardBody>
      </Card>
      {year && pcrs.data && (
        <Box gap={'medium'} pad={'medium'} align={'start'} fill={'horizontal'}>
          <InfantPcrsList pcrs={pcrs} />
        </Box>
      )}
    </Box>
  );
};

export default InfantPcrs;
