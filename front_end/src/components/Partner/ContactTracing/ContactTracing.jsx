import {
  Box,
  Button,
  CardBody,
  CardHeader,
  Heading,
  Layer,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import AppCard from '../../AppCard/AppCard';
import Layout from '../../Layout/Layout';
import Spinner from '../../Spinner';
import PartnerTabs from '../PartnerTabs/PartnerTabs';
import { Add, Close, Edit } from 'grommet-icons';
import ContactTracingEdit from './ContactTracingEdit';
import { parseDate } from '../../../dates';

const row = (data, onClickEdit) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text size={'small'} align={'start'}>
          {data.test}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'} align={'start'}>
          {data.testResult}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'} align={'start'}>
          {data.comments}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'} align={'start'}>
          {data.date ? parseDate(data.date, 'dd LLL yyyy') : 'N/A'}
        </Text>
      </TableCell>
      <TableCell align={'center'} onClick={() => onClickEdit(data)}>
        <Edit />
      </TableCell>
    </TableRow>
  );
};

const ContactTracingTable = ({ children, data, onClickEdit }) => {
  if (!data) {
    return (
      <Box alignContent={'center'}>
        <Text>No data found.</Text>
      </Box>
    );
  }

  return (
    <Box gap={'medium'} pad={'medium'} align={'start'} fill>
      {children}
      <Table>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text align={'start'}>Test</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Result</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Comments</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Date</Text>
            </TableCell>
            <TableCell align={'start'}></TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{data.map((d) => row(d, onClickEdit))}</TableBody>
      </Table>
    </Box>
  );
};

const ContactTracing = () => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [tracings, setTracings] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });
  const [editingTracing, setEditingTracing] = React.useState(undefined);
  const history = useHistory();

  const onClickEdit = (tracing) => setEditingTracing(tracing);
  const onCloseEditScreen = () => {
    setEditingTracing(undefined);
    setTracings({ data: undefined, loading: true, error: undefined });
  };

  React.useEffect(() => {
    const getTracings = () => {
      httpInstance
        .get(`/partners/${patientId}/contactTracing`)
        .then((resp) => {
          setTracings({
            data: resp.data,
            loading: false,
            error: undefined,
          });
        })
        .catch((e) => {
          console.error(e);
          setTracings({ data: undefined, loading: false, error: e.toJSON() });
        });
    };
    if (tracings.loading) {
      getTracings();
    }
  }, [httpInstance, patientId, tracings]);

  if (tracings.loading) {
    return (
      <Layout>
        <Box
          direction={'column'}
          gap={'large'}
          pad={'large'}
          justify={'center'}
          align={'center'}
          fill
        >
          <Heading>
            <Text>Loading....</Text>
            <Spinner />
          </Heading>
        </Box>
      </Layout>
    );
  }

  if (tracings.error) {
    return (
      <Box
        direction={'colomn'}
        gap={'large'}
        pad={'large'}
        justify={'center'}
        align={'center'}
        fill
      >
        <Heading>
          <Text>Ooops. An error occurred while loading the data. </Text>
        </Heading>
      </Box>
    );
  }

  return (
    <Layout>
      <Box
        direction={'column'}
        gap={'xxlarge'}
        pad={{ left: 'small', bottom: 'xxsmall' }}
        justify={'evenly'}
        align={'center'}
        fill
      >
        <PartnerTabs data={tracings.data}>
          <AppCard overflow={'scroll'} pad={'small'} fill={'horizontal'}>
            <CardHeader justify={'start'} pad={'medium'}>
              <Box direction={'row'} align={'start'} fill={'horizontal'}>
                <Box fill={'horizontal'}>
                  <span>
                    <Text size={'xxlarge'} weight={'bold'}>
                      Partner Contact Tracing
                    </Text>
                  </span>
                  <Text size={'large'}>
                    {' '}
                    {tracings.data.patient.firstName}{' '}
                    {tracings.data.patient.lastName}
                  </Text>
                </Box>
                <Box
                  align={'start'}
                  fill={'horizontal'}
                  direction={'row-reverse'}
                >
                  <Button
                    icon={<Add />}
                    label={'Add'}
                    onClick={() =>
                      history.push(
                        `/patient/${patientId}/partners/contactTracing/new`
                      )
                    }
                  />
                </Box>
              </Box>
            </CardHeader>

            <CardBody gap={'medium'} pad={'medium'}>
              {editingTracing && (
                <Layer
                  position={'right'}
                  full={'vertical'}
                  onClickOutside={onCloseEditScreen}
                  modal
                >
                  <Box
                    fill={'vertical'}
                    overflow={'auto'}
                    width={'medium'}
                    pad={'medium'}
                  >
                    <Box
                      flex={false}
                      direction={'row-responsive'}
                      justify={'between'}
                    >
                      <Heading level={2} margin={'none'}>
                        Edit
                      </Heading>
                      <Button icon={<Close />} onClick={onCloseEditScreen} />
                    </Box>
                    <ContactTracingEdit
                      contactTracing={editingTracing}
                      closeEditScreen={onCloseEditScreen}
                    />
                  </Box>
                </Layer>
              )}
              <ContactTracingTable
                data={tracings.data.contactTracing}
                onClickEdit={onClickEdit}
              />
            </CardBody>
          </AppCard>
        </PartnerTabs>
      </Box>
    </Layout>
  );
};

export default ContactTracing;
