import React from 'react';
import {
  Box,
  CardBody,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import format from 'date-fns/format';
import parseISO from 'date-fns/parseISO';
import AppCard from '../../../AppCard/AppCard';
import { View } from 'grommet-icons';
import { useHistory } from 'react-router-dom';

const row = (data, onClick) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.lmp), 'dd LLL yyyy')}</Text>
      </TableCell>
      <TableCell>
        <Text>
          {data.edd ? format(parseISO(data.edd), 'dd LLL yyyy') : 'N/A'}
        </Text>
      </TableCell>
      <TableCell onClick={() => onClick(data)}>
        <View />
      </TableCell>
    </TableRow>
  );
};

const PregnanciesTable = ({ children, data, onClick }) => {
  return (
    <Box gap={'medium'} align={'center'}>
      {children}
      <Table caption={'Pregnancies'}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>Last Menstrual Period</TableCell>
            <TableCell align={'start'}>Estimated Time of Delivery</TableCell>
            <TableCell />
          </TableRow>
        </TableHeader>
        <TableBody>{data.map((d) => row(d, onClick))}</TableBody>
      </Table>
    </Box>
  );
};

const PregnancyHistory = (props) => {
  const { pregnancies } = props;
  const history = useHistory();

  const onClick = (d) =>
    history.push(`/patient/${d.patientId}/pregnancy/${d.id}`);

  if (pregnancies.length === 0) {
    return <>No Pregnancies</>;
  }

  return (
    <AppCard fill={'horizontal'}>
      <CardBody gap={'medium'} pad={'medium'}>
        <PregnanciesTable data={pregnancies} onClick={onClick} />
      </CardBody>
    </AppCard>
  );
};

export default PregnancyHistory;
