import React from 'react';
import { useHttpApi } from '../../providers/HttpProvider';
import Spinner from '../Spinner';
import {
  Box,
  Button,
  Heading,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
  Layer,
} from 'grommet';
import { Close, Edit } from 'grommet-icons';
import UserEdit from './UserEdit';

const UserListData = ({ children }) => {
  const [users, setUsers] = React.useState([]);
  const [error, setError] = React.useState();
  const [loading, setLoading] = React.useState(true);
  const { httpInstance } = useHttpApi();
  React.useEffect(() => {
    const getUsers = () => {
      httpInstance
        .get('/admin/users')
        .then((result) => {
          setUsers(result.data);
          setLoading(false);
          setError(undefined);
        })
        .catch((error) => {
          setError(error);
          setUsers([]);
          setLoading(false);
        });
    };
    if (loading) {
      getUsers();
    }
  }, [httpInstance, users, loading]);
  return children({ users, loading, error });
};

const userRow = (data, onClickEdit) => {
  return (
    <TableRow key={data.id}>
      <TableCell size={'medium'}>
        <Text size={'small'}>{data.email}</Text>
      </TableCell>
      <TableCell size={'small'}>
        <Text size={'small'}>{data.firstName}</Text>
      </TableCell>
      <TableCell size={'small'}>
        <Text size={'small'}>{data.lastName}</Text>
      </TableCell>
      <TableCell size={'xxsmall'} onClick={() => onClickEdit(data)}>
        <Edit />
      </TableCell>
    </TableRow>
  );
};

const UserTable = ({ users, onClickEdit }) => {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableCell size={'medium'}>
            <Text weight={'bold'} size={'small'}>
              Email
            </Text>
          </TableCell>
          <TableCell size={'medium'}>
            <Text weight={'bold'} size={'small'}>
              First Name
            </Text>
          </TableCell>
          <TableCell size={'medium'}>
            <Text weight={'bold'} size={'small'}>
              Last Name
            </Text>
          </TableCell>
          <TableCell size={'xxsmall'} />
        </TableRow>
      </TableHeader>
      <TableBody>{users.map((u) => userRow(u, onClickEdit))}</TableBody>
    </Table>
  );
};

const Loading = () => {
  return (
    <Box
      align={'center'}
      justify={'center'}
      direction={'column'}
      gap={'large'}
      pad={'large'}
      fill={'horizontal'}
    >
      <Spinner />
    </Box>
  );
};

const UserListComponent = ({ users, loading, error }) => {
  const [editingUser, setEditingUser] = React.useState();
  const onClickEdit = (user) => setEditingUser(user);
  return (
    <>
      {loading && <Loading />}
      {!loading && error && <Text>Error!</Text>}
      {!loading && !error && (
        <Box
          direction={'column'}
          gap={'xxsmall'}
          pad={'xxsmall'}
          fill={'vertical'}
        >
          <Box
            align={'start'}
            justify={'start'}
            gap={'xxsmall'}
            pad={'xxsmall'}
            direction={'row-responsive'}
            margin={{ left: 'small', top: 'xxsmall' }}
          >
            <Heading level={3} responsive={true}>
              <Text size={'xxlarge'}>Users</Text>
            </Heading>
          </Box>
          <Box
            align={'start'}
            gap={'xxsmall'}
            pad={'xxsmall'}
            margin={{ left: 'small', bottom: 'medium' }}
          >
            <Button secondary label={'Add User'} type={'button'} />
          </Box>
          <Box
            align={'center'}
            justify={'center'}
            gap={'medium'}
            pad={'medium'}
            margin={{ left: 'small', bottom: 'medium' }}
            fill={'horizontal'}
          >
            {editingUser && (
              <Layer
                position={'right'}
                full={'vertical'}
                onClickOutside={() => setEditingUser(undefined)}
                modal
              >
                <Box
                  as={'form'}
                  fill={'vertical'}
                  overflow={'auto'}
                  width={'medium'}
                  pad={'medium'}
                >
                  <Box flex={false} direction={'row'} justify={'between'}>
                    <Heading level={2} margin={'none'}>
                      Edit
                    </Heading>
                    <Button
                      icon={<Close />}
                      onClick={() => setEditingUser(undefined)}
                    />
                  </Box>
                  <UserEdit user={editingUser} />
                </Box>
              </Layer>
            )}
            <UserTable users={users} onClickEdit={onClickEdit} />
          </Box>
        </Box>
      )}
    </>
  );
};

const UserList = () => {
  return (
    <UserListData>
      {(renderProps) => <UserListComponent {...renderProps} />}
    </UserListData>
  );
};

export default UserList;
