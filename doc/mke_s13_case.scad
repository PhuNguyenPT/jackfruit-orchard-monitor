// =====================================================
// MKE-S13 Capacitive Soil Moisture Sensor — Case
// Parametric OpenSCAD  |  v5.1 (Drop-In U-Slot & Pillars)
// Units: mm
// =====================================================

include <mke_s13_config.scad>

$fn = 48;

// =====================================================
// 2D PROFILES
// =====================================================
module pcb_2d() {
    hull() {
        translate([0,        pcb_w/2])         circle(r=0.1);
        translate([chev_l,   0])                circle(r=0.1);
        translate([chev_l,   pcb_w])           circle(r=0.1);
        // X-axis extended to ensure rapid-fab PCBs don't crash into the wall
        translate([pcb_l - corner_r + fab_x_tol, corner_r])           circle(r=corner_r);
        translate([pcb_l - corner_r + fab_x_tol, pcb_w - corner_r])   circle(r=corner_r);
    }
}

module pcb_box_section_2d() {
    intersection() {
        pcb_2d();
        translate([safe_line_x - 1, -5])
            square([box_l + 6, pcb_w + 10]);
    }
}

module shell_2d(w) {
    offset(r = w + pcb_gap) pcb_box_section_2d();
}

// =====================================================
// BOTTOM SHELL (Clean exterior, Drop-in U-Slot)
// =====================================================
module bottom_shell() {
    difference() {
        // Main Outer body box
        linear_extrude(outer_h) shell_2d(wall);

        // Interior hollow
        translate([0, 0, floor_t])
            linear_extrude(outer_h)
                offset(r = pcb_gap) pcb_box_section_2d();

        // --- DROP-IN U-SLOT FOR PCB ---
        // Open to the top rim to allow Z-axis insertion
        translate([
            safe_line_x - 1 - pcb_gap - wall - 1.0,
            pcb_w/2 - (pcb_w/2 + slot_gap),
            floor_t + pcb_t + 0.2 - slot_gap
        ])
            cube([
                wall + 2.0,
                pcb_w + 2*slot_gap,
                outer_h
            ]);

        // NOTE: The "Pin alignment holes through the floor" block
        // was removed here to make the bottom shell 100% waterproof!

        // Alignment nub dimples on rim
        translate([nub_x, -wall*0.5 - pcb_gap, outer_h - nub_h]) cylinder(d = nub_d + nub_clearance, h = nub_h + 0.1);
        translate([nub_x, pcb_w + wall*0.5 + pcb_gap, outer_h - nub_h]) cylinder(d = nub_d + nub_clearance, h = nub_h + 0.1);
    }

    // PCB mounting standoffs (Pins from the lid will seat inside here)
    // Updated to dynamically calculate height based on z_pcb_seat
    // to properly include the solder_z_tol buffer.
    standoff_h = z_pcb_seat - floor_t;

    translate([0, 0, floor_t]) {
        for (y_off = [-hole_sp/2, hole_sp/2]) {
            translate([hole_x, hole_cy + y_off, 0])
            difference() {
                cylinder(d = pcb_boss_d, h = standoff_h);
                // The bore hole for the lid pin now stops at the solid floor (blind hole)
                translate([0, 0, -0.01])
                    cylinder(d = hole_d, h = standoff_h + 1);
            }
        }
    }

    // --- CONNECTOR / PCB PARTITION BULKHEAD (LOWER HALF) ---
    // Matching lower half of the bulkhead in lid() -- rooted in the floor
    // (z=0) and rising up to exactly z_pcb_seat, the SAME level the PCB
    // mounting standoffs above present as their top face. This is a flush,
    // zero-clearance support surface (not a sealing gap) so the PCB rests
    // level across both the standoffs and this wall, with no rocking or
    // height mismatch. The lid's upper half of the bulkhead meets this
    // piece's top at the PCB's underside, with the PCB itself (thickness
    // pcb_t) plus a slot_gap clearance occupying the open band in between
    // (see lid()'s partition_drop_h calculation).
    intersection() {
        translate([partition_x1, -1000, 0])
            cube([partition_t, 2000, z_pcb_seat]);
        linear_extrude(z_pcb_seat)
            offset(r = pcb_gap) pcb_box_section_2d();
    }
}

// =====================================================
// LID (Locking pillars, Rotated Slot, Closure Baffle)
// =====================================================
module lid() {
    difference() {
        union() {
            // Main Lid plate
            linear_extrude(lid_t) shell_2d(wall);

            // --- CLOSURE BAFFLE ---
            // Drops down into the bottom shell's U-slot to seal the case
            baffle_h = inner_h - (2 * pcb_t + 0.2 + slot_gap);
            baffle_w = pcb_w + 2*slot_gap - baffle_clearance; // baffle_clearance mm tolerance total for smooth fit

            translate([
                safe_line_x - 1 - pcb_gap - wall,
                pcb_w/2 - baffle_w/2,
                -baffle_h
            ])
                cube([wall, baffle_w, baffle_h]);

            // --- CONNECTOR / PCB PARTITION BULKHEAD (UPPER HALF) ---
            // Splash/dust barrier between the connector cavity (x > partition_x2,
            // open to outside air via the ROTATED CONNECTOR SLOT below) and the
            // main PCB cavity (x < partition_x1). The PCB + male connector are
            // pre-soldered into one rigid unit before assembly, so the bottom
            // shell must stay a clean unobstructed box for that unit to drop
            // straight into -- so this wall only covers the space ABOVE the
            // PCB, hanging from the lid down to just shy of the board's top
            // face (slot_gap clearance, purely a sealing gap, not load-bearing).
            // It never needs to clear the connector either, since the connector
            // sits entirely at x > partition_x2, outside this wall's footprint.
            // The matching LOWER half (floor up to the PCB's underside) is a
            // separate piece rooted in bottom_shell() -- see pcb_seat level
            // there, which this wall's bottom edge is calculated to meet
            // exactly flush, with zero gap, so the PCB has continuous level
            // support across the standoffs and this bulkhead alike.
            partition_drop_h = outer_h - (z_pcb_seat + pcb_t + slot_gap);
            partition_w      = pcb_w + 2*pcb_gap - baffle_clearance;

            translate([
                partition_x1,
                pcb_w/2 - partition_w/2,
                -partition_drop_h
            ])
                cube([partition_t, partition_w, partition_drop_h]);
        }

        // ROTATED CONNECTOR SLOT (Wide along Y axis to fit the plug)
        // Inner (tip-side) edge is pinned flush to partition_x2 -- the
        // connector-side face of the new bulkhead -- so this opening only
        // ever exposes the connector cavity, never the sealed main cavity.
        translate([
            partition_x2,
            pcb_w/2 - (conn_l / 2) - cable_clear,
            -0.01
        ])
            cube([
                (pcb_l + cable_clear) - partition_x2,
                conn_l + 2*cable_clear,
                lid_t + 0.1
            ]);
    }

    // Alignment registration nubs
    translate([nub_x, -wall*0.5 - pcb_gap, -nub_h]) cylinder(d = nub_d, h = nub_h);
    translate([nub_x, pcb_w + wall*0.5 + pcb_gap, -nub_h]) cylinder(d = nub_d, h = nub_h);

    // --- INTEGRATED LOCKING PILLARS ---
    // Protrudes downwards from the underside of the lid (Z=0 local)
    for (y_off = [-hole_sp/2, hole_sp/2]) {
        translate([hole_x, hole_cy + y_off, 0]) {

            // 1. Upper Wide Shoulder: Clamps onto the top face of the PCB
            // Dynamically calculated to drop exactly to the board surface
            shoulder_h = inner_h - (z_pcb_seat + pcb_t);
            translate([0, 0, -shoulder_h])
                cylinder(d = pcb_boss_d, h = shoulder_h);

            // 2. Interlocking Pin: Passes through PCB and bottom standoffs
            // Drawn from the bottom UP. Bottom is lifted by squish_tol.
            pin_h = inner_h - shoulder_h - squish_tol;
            translate([0, 0, -inner_h + squish_tol])
                cylinder(d = lock_pin_d, h = pin_h);
        }
    }
}

// =====================================================
// PRINT-READY ORIENTATIONS
// =====================================================
// The lid() module above is authored in "assembly logic":
// the flat plate sits at the TOP (z=0..lid_t) and the
// closure baffle + locking pillars hang DOWNWARD into
// negative Z, matching how the lid sits when placed onto
// the bottom shell. That orientation is NOT printable as-is
// (parts fall below the z=0 bed plane), and printing it
// pillars-down would require support material under every
// pillar tip.
//
// lid_print_ready() mirrors the part so the plate becomes
// the base sitting flat on the bed, and the pillars become
// self-supporting upward-growing towers — no supports needed.
// Print the plate first; the pillars build up naturally on
// top of it, same logic as the bottom shell's standoffs.
module lid_print_ready() {
    translate([0, 0, lid_t])
        mirror([0, 0, 1])
            lid();
}

// =====================================================
// RENDER
// =====================================================
color("SteelBlue",  0.85) bottom_shell();
color("LightBlue",  0.70) translate([0, pcb_w + 15, 0]) lid_print_ready();