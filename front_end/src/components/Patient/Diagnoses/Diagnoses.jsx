import {
  Box,
  Card,
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

const DiagnosisRow = ({ data }) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.date), 'dd LLL yyyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <TableCell>
          <Text align={'start'}>{data.name}</Text>
        </TableCell>
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
              <Text>Date</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Illness</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data.map((d) => (
            <DiagnosisRow data={d} />
          ))}
        </TableBody>
      </Table>
    </Box>
  );
};

const DiagnosisHistory = (props) => {
  const { diagnosisHistory, caption } = props;
  return (
    <Card>
      <CardBody gap={'medium'} pad={'medium'}>
        <DiagnosisTable data={diagnosisHistory} caption={caption}>
          <CircleInformation size={'large'} />
        </DiagnosisTable>
      </CardBody>
    </Card>
  );
};

export default DiagnosisHistory;
