import {
  Box,
  CardBody,
  Table,
  TableCell,
  TableHeader,
  TableRow,
  TableBody,
  Text,
} from 'grommet';
import { CircleInformation } from 'grommet-icons';
import React from 'react';
import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';
import AppCard from '../../AppCard/AppCard';

const diagnosisRow = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text size={'small'}>{format(parseISO(data.date), 'dd LLL yyyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text align={'start'} size={'small'}>
          {data.name}
        </Text>
      </TableCell>
    </TableRow>
  );
};

const DiagnosisTable = ({ children, data, caption, ...rest }) => {
  return (
    <Box gap={'medium'} align={'center'} {...rest}>
      {children}
      <Table caption={caption}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text size={'small'}>Date</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Illness</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{data.map((d) => diagnosisRow(d))}</TableBody>
      </Table>
    </Box>
  );
};

const DiagnosisHistory = (props) => {
  const { diagnosisHistory, caption } = props;

  return (
    <AppCard fill={'horizontal'}>
      <CardBody gap={'medium'} pad={'medium'}>
        <DiagnosisTable data={diagnosisHistory} caption={caption}>
          <CircleInformation size={'large'} />
        </DiagnosisTable>
      </CardBody>
    </AppCard>
  );
};

export default DiagnosisHistory;
