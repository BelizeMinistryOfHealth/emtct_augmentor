import React from 'react';
import { Anchor, Box, Header, Nav } from 'grommet';
import { withAuth0 } from '@auth0/auth0-react';
import LogoutButton from '../Logout/Logout';

const items = [
    {label: 'HIV', href: '#'},
    {label: 'Syphillis', href: '#'},
];


const Navbar = (props) => {
    const { user, isAuthenticated } = props.auth0;

    return isAuthenticated && (
    <Header background={'dark-1'} pad={'small'}>
        <Box direction={'row'} align={'center'} gap={'small'}>
            <Anchor color={'white'}>
                EMTCT Project | {user.name}
            </Anchor>
        </Box>
        <Nav direction={'row'} align={'baseline'} gap={'small'}>
            {items.map(item =>(
                <Anchor href={item.href} label={item.label} key={item.label} />
            ))}
        </Nav>
        <Box direction={'row'} align={'end'} gap={'small'}>
            <LogoutButton />
        </Box>
        
    </Header>);
};

export default withAuth0(Navbar);