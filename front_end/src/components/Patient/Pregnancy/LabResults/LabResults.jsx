import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';
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

const labResultsRow = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text align={'start'}>
          {format(parseISO(data.dateSampleTaken), 'dd LLL yyyy')}
        </Text>
      </TableCell>
      <TableCell>
        <Text align={'start'}>{data.testResult}</Text>
      </TableCell>
      <TableCell>
        <Text align={'start'}>{data.testName}</Text>
      </TableCell>
    </TableRow>
  );
};

const LabResultsTable = ({ children, data, caption, ...rest }) => {
  console.dir({ data });
  return (
    <Box gap={'medium'} pad={'medium'} align={'center'} {...rest}>
      {children}
      <Table caption={caption}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text align={'start'}> Date sample taken</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}> Test Result</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Test Name</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{data.map((d) => labResultsRow(d))}</TableBody>
      </Table>
    </Box>
  );
};

const LabResults = (props) => {
  const { labResults, caption } = props;
  return (
    <Card>
      <CardBody gap={'medium'} pad={'medium'}>
        <LabResultsTable data={labResults} caption={caption} />
      </CardBody>
    </Card>
  );
};

export default LabResults;
