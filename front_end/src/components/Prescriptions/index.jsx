import { format, parseISO } from 'date-fns';
import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import React from 'react';

const row = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.pharmaceutical}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.strength}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>
          {format(parseISO(data.prescribedTime), 'dd LLL yyy')}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.comments}</Text>
      </TableCell>
    </TableRow>
  );
};

const Prescriptions = ({ children, data, ...rest }) => {
  const prescriptions = data.prescriptions ?? [];

  return (
    <Box gap={'medium'} align={'center'} {...rest} fill={'horizontal'}>
      {children}
      <Table>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text>Treatment</Text>
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
        <TableBody>{prescriptions.map((d) => row(d))}</TableBody>
      </Table>
    </Box>
  );
};

export default Prescriptions;
