import React from 'react';
import Spinner from '../../Spinner';
import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import { parseISO } from 'date-fns';
import format from 'date-fns/format';
import { EmptyCircle } from 'grommet-icons';

const row = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell>
        <Text size={'small'}>{data.motherName}</Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>{data.infantName}</Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>
          {data.PCR1DateSampleTaken ? (
            format(parseISO(data.PCR1DateSampleTaken), 'dd LLL yyyy')
          ) : (
            <EmptyCircle color={'red'} />
          )}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>
          {data.PCR1DueDate
            ? format(parseISO(data.PCR1DueDate), 'dd LLL yyyy')
            : 'Not due this year'}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>
          {data.PCR2DateSampleTaken ? (
            format(parseISO(data.PCR2DateSampleTaken), 'dd LLL yyyy')
          ) : (
            <EmptyCircle color={'red'} />
          )}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>
          {data.PCR2DueDate
            ? format(parseISO(data.PCR2DueDate), 'dd LLL yyyy')
            : 'Not due this year'}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>
          {data.PCR3DateSampleTaken ? (
            format(parseISO(data.PCR3DateSampleTaken), 'dd LLL yyyy')
          ) : (
            <EmptyCircle color={'red'} />
          )}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>
          {data.PCR3DueDate
            ? format(parseISO(data.PCR3DueDate), 'dd LLL yyyy')
            : 'Not due this year'}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>
          {data.ELISADateSampleTaken ? (
            format(parseISO(data.ELISADateSampleTaken), 'dd LLL yyyy')
          ) : (
            <EmptyCircle color={'red'} />
          )}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>
          {data.ELISADueDate
            ? format(parseISO(data.ELISADueDate), 'dd LLL yyyy')
            : 'Not due this year'}
        </Text>
      </TableCell>
    </TableRow>
  );
};

const PcrTable = ({ data }) => {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableCell>
            <Text size={'small'}>Mother Name</Text>
          </TableCell>
          <TableCell>
            <Text size={'small'}>Child Name</Text>
          </TableCell>
          <TableCell>
            <Text size={'small'}>PCR 1: Sample Taken</Text>
          </TableCell>
          <TableCell>
            <Text size={'small'}>PCR 1: Due Date</Text>
          </TableCell>
          <TableCell size={'small'}>
            <Text size={'small'}>PCR 2: Sample Taken</Text>
          </TableCell>
          <TableCell>
            <Text size={'small'}>PCR 2: Due Date</Text>
          </TableCell>
          <TableCell>
            <Text size={'small'}>PCR 3: Sample Taken</Text>
          </TableCell>
          <TableCell>
            <Text size={'small'}>PCR 3: Due Date</Text>
          </TableCell>
          <TableCell>
            <Text size={'small'}>ELISA: Sample Taken</Text>
          </TableCell>
          <TableCell>
            <Text size={'small'}>ELISA: Due Date</Text>
          </TableCell>
        </TableRow>
      </TableHeader>
      {data && data.length > 0 && (
        <TableBody>{data.map((d) => row(d))}</TableBody>
      )}
    </Table>
  );
};

/**
 * InfantPcrsList displays the HIV Screenings due for a given year.
 * @param pcrs
 * @returns {JSX.Element}
 * @constructor
 */
const InfantPcrsList = ({ pcrs }) => {
  const { loading, error } = pcrs;
  return (
    <>
      {loading && (
        <Box align={'center'} fill={'horizontal'}>
          <Spinner />
        </Box>
      )}
      {!loading && !error && pcrs.data.length === 0 && (
        <Box align={'center'} fill={'horizontal'}>
          <Text>No Results Found.</Text>
        </Box>
      )}
      {!loading && !error && pcrs.data.length > 0 && (
        <>
          <PcrTable data={pcrs.data} />
        </>
      )}
    </>
  );
};

export default InfantPcrsList;
