import { format, parseISO } from 'date-fns';
import {
  Box,
  Card,
  CardBody,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import React from 'react';
import { useRecoilValueLoadable } from 'recoil';
import { arvTreatmentsSelector } from '../../../state';

const row = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text>{data.arvName}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.dosage}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.date), 'dd LLL yyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.comments}</Text>
      </TableCell>
    </TableRow>
  );
};

const ArvsTable = ({ children, data, ...rest }) => {
  return (
    <Box gap={'medium'} align={'center'} {...rest}>
      {children}
      <Table caption={'Arv Treatment'}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text>ARV</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Dosage</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Date</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Comments</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{data.map((d) => row(d))}</TableBody>
      </Table>
    </Box>
  );
};

const ArvTreatment = (props) => {
  const { patientId, encounterId } = props;
  const { state, contents } = useRecoilValueLoadable(
    arvTreatmentsSelector(patientId, encounterId)
  );

  switch (state) {
    case 'hasValue':
      return (
        <Card>
          <CardBody gap={'medium'} pad={'medium'}>
            <ArvsTable data={contents}></ArvsTable>
          </CardBody>
        </Card>
      );
    case 'loading':
      return 'loading';
    case 'hasError':
      return contents.message;
    default:
      return '';
  }
};

export default ArvTreatment;
