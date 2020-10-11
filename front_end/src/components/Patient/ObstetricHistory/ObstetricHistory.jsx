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
import { History } from 'grommet-icons';
import React from 'react';
import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';

const eventRow = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.date), 'dd LLL yyyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.event}</Text>
      </TableCell>
    </TableRow>
  );
};

const EventTable = ({ children, data, ...rest }) => {
  return (
    <Box gap={'medium'} align={'center'} {...rest}>
      {children}
      <Table caption={'Obstetric History'}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text>Date</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Obstetric Event</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{data.map((d) => eventRow(d))}</TableBody>
      </Table>
    </Box>
  );
};

const ObstetricHistory = (props) => {
  const { obstetricHistory } = props;
  return (
    <Card>
      <CardBody gap={'medium'} pad={'medium'}>
        <EventTable data={obstetricHistory}>
          <History size={'large'} />
        </EventTable>
      </CardBody>
    </Card>
  );
};

export default ObstetricHistory;
